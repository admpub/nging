package gerberos

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ipMagicText = "%ip%"
	idMagicText = "%id%"
)

var (
	ipMagicRegexp = regexp.MustCompile(ipMagicText)
	ipRegexpText  = `(?P<ip>(\d?\d?\d\.){3}\d?\d?\d|\[?([0-9A-Fa-f]{0,4}::?){1,6}[0-9A-Fa-f]{0,4}::?[0-9A-Fa-f]{0,4})\]?`
	idMagicRegexp = regexp.MustCompile(idMagicText)
	idRegexpText  = `(?P<id>(.*))`
)

type Rule struct {
	Source      []string
	Regexp      []string
	Action      []string
	Aggregate   []string
	Occurrences []string

	runner      *Runner
	name        string
	source      Source
	regexp      []*regexp.Regexp
	action      Action
	aggregate   *aggregate
	occurrences *occurrences
}

func (r *Rule) initializeSource() error {
	if r.Source == nil {
		return ErrMissingSource
	}

	if len(r.Source) == 0 {
		return ErrEmptySource
	}
	sfn, ok := sources[r.Source[0]]
	if !ok {
		return fmt.Errorf(`%w: %v`, ErrUnknownSource, r.Source[0])
	}
	r.source = sfn()
	return r.source.Initialize(r)
}

func (r *Rule) initializeRegexp() error {
	if r.Regexp == nil {
		return ErrMissingRegexp
	}

	if len(r.Regexp) == 0 {
		return ErrEmptyRegexp
	}

	r.regexp = make([]*regexp.Regexp, 0, len(r.Regexp))
	for _, s := range r.Regexp {
		if strings.Contains(s, "(?P<ip>") {
			return errors.New(`regexp must not contain a subexpression named "ip" ("(?P<ip>")`)
		}

		if strings.Contains(s, "(?P<id>") {
			return errors.New(`regexp must not contain a subexpression named "id" ("(?P<id>")`)
		}

		if len(ipMagicRegexp.FindAllStringIndex(s, -1)) != 1 {
			return fmt.Errorf(`"%s" must appear exactly once in regexp`, ipMagicText)
		}

		if r.Aggregate != nil && len(idMagicRegexp.FindAllStringIndex(s, -1)) != 1 {
			return fmt.Errorf(`"%s" must appear exactly once in regexp if the aggregate option is used`, idMagicText)
		}

		t := strings.Replace(s, ipMagicText, ipRegexpText, 1)
		t = strings.Replace(t, idMagicText, idRegexpText, 1)
		re, err := regexp.Compile(t)
		if err != nil {
			return err
		}
		r.regexp = append(r.regexp, re)
	}

	return nil
}

func (r *Rule) initializeAction() error {
	if r.Action == nil {
		return ErrMissingAction
	}

	if len(r.Action) == 0 {
		return ErrEmptyAction
	}
	afn, ok := actions[r.Action[0]]
	if !ok {
		return fmt.Errorf(`%w: %v`, ErrUnknownAction, r.Action[0])
	}
	r.action = afn()
	return r.action.Initialize(r)
}

func (r *Rule) initializeAggregate() error {
	if r.Aggregate == nil {
		return nil
	}

	if len(r.Aggregate) < 1 {
		return ErrMissingIntervalParameter
	}
	i, err := time.ParseDuration(r.Aggregate[0])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidIntervalParameter, err)
	}

	if len(r.Aggregate) < 2 {
		return ErrMissingRegexp
	}

	res := make([]*regexp.Regexp, 0, len(r.Aggregate)-1)
	for _, s := range r.Aggregate[1:] {
		if strings.Contains(s, "(?P<id>") {
			return errors.New(`regexp must not contain a subexpression named "id" ("(?P<id>")`)
		}

		if len(idMagicRegexp.FindAllStringIndex(s, -1)) != 1 {
			return fmt.Errorf(`"%s" must appear exactly once in regexp`, idMagicRegexp)
		}

		re, err := regexp.Compile(strings.Replace(s, idMagicText, idRegexpText, 1))
		if err != nil {
			return err
		}
		res = append(res, re)
	}

	r.aggregate = newAggregate(i, res)

	return nil
}

func (r *Rule) initializeOccurrences() error {
	if r.Occurrences == nil {
		return nil
	}

	if len(r.Occurrences) < 1 {
		return ErrMissingCountParameter
	}
	c, err := strconv.Atoi(r.Occurrences[0])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidCountParameter, err)
	}
	if c < 2 {
		return fmt.Errorf("%w: must be > 1", ErrInvalidCountParameter)
	}

	if len(r.Occurrences) < 2 {
		return ErrMissingIntervalParameter
	}
	i, err := time.ParseDuration(r.Occurrences[1])
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidIntervalParameter, err)
	}

	r.occurrences = newOccurrences(i, c)

	return nil
}

func (r *Rule) initialize(rn *Runner) error {
	r.runner = rn

	if err := r.initializeSource(); err != nil {
		return err
	}

	if err := r.initializeRegexp(); err != nil {
		return err
	}

	if err := r.initializeAction(); err != nil {
		return err
	}

	if err := r.initializeAggregate(); err != nil {
		return err
	}

	if err := r.initializeOccurrences(); err != nil {
		return err
	}

	return nil
}

func (r *Rule) ProcessScanner(name string, args ...string) (chan *Match, error) {
	stop := make(chan bool, 1)

	cmd := exec.Command(name, args...)
	o, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	e, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}
	log.Printf(`%s: scanning process stdout and stderr: "%s"`, r.name, cmd)

	go func() {
		select {
		case <-stop:
		case <-r.runner.stopped.Done():
		}
		if cmd.Process != nil {
			cmd.Process.Signal(os.Interrupt)
			time.Sleep(5 * time.Second)
			select {
			case <-stop:
			default:
				cmd.Process.Kill()
			}
		}
	}()

	c := make(chan *Match, 1)
	go func() {
		defer func() {
			stop <- true
			close(stop)
		}()

		sc := bufio.NewScanner(o)
		for sc.Scan() {
			if m, err := r.Match(sc.Text()); err == nil {
				c <- m
			} else {
				if r.runner.Configuration.Verbose {
					log.Printf("%s: failed to create match: %s", r.name, err)
				}
			}
		}
		close(c)
		if err = sc.Err(); err != nil {
			log.Printf(`%s: error while scanning command "%s": %s`, r.name, cmd, err.Error())
		}
		if err = cmd.Wait(); err != nil {
			var eerr *exec.ExitError
			if errors.As(err, &eerr) {
				if eerr.ProcessState.ExitCode() == -1 {
					// The process was terminated by a signal. This is part of a graceful
					// shutdown. Therefore it isn't logged.
					return
				}
			}
			log.Printf(`%s: error while executing command "%s": %s`, r.name, cmd, err.Error())
		}
	}()
	go func() {
		sc := bufio.NewScanner(e)
		for sc.Scan() {
			log.Printf(`%s: process stderr: "%s"`, r.name, sc.Text())
		}
	}()

	return c, nil
}

func (r *Rule) worker(requeue bool) error {
	c, err := r.source.Matches()
	if err != nil {
		log.Printf("%s: failed to initialize matches channel: %s", r.name, err)
		return err
	}

	for m := range c {
		p := true
		if r.occurrences != nil {
			p = r.occurrences.add(m.IP)
		}

		if p {
			if err := r.action.Perform(m); err != nil {
				log.Printf("%s: failed to perform action: %s", r.name, err)
			}
		}
	}

	if requeue {
		log.Printf("%s: queuing worker for respawn", r.name)
		select {
		case r.runner.respawnWorkerChan <- r:
		case <-r.runner.stopped.Done():
		}
	}

	return nil
}
