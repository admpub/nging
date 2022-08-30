package redis

import "github.com/gomodule/redigo/redis"

// ListChannels 获取频道列表
func (r *Redis) ListChannels(pattern ...interface{}) ([]string, error) {
	reply, err := redis.Strings(r.conn.Do("PUBSUB", append([]interface{}{`CHANNELS`}, pattern...)...))
	if err != nil {
		return nil, err
	}
	return reply, err
}

// ChannelSubscribeNum 获取频道订阅量
func (r *Redis) ChannelSubscribeNum(channels ...interface{}) (map[string]int64, error) {
	reply, err := redis.Int64Map(r.conn.Do("PUBSUB", append([]interface{}{`NUMSUB`}, channels...)...))
	if err != nil {
		return nil, err
	}
	return reply, err
}
