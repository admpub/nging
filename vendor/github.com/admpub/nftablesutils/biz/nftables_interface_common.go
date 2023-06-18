package biz

import (
	"fmt"

	"github.com/google/nftables"
)

func (nft *NFTables) applyCommonRules(c *nftables.Conn, iface string) error {
	err := nft.inputHostBaseRules(c, nft.wanIface)
	if err != nil {
		return fmt.Errorf(`nft.inputHostBaseRules(%q): %w`, nft.wanIface, err)
	}
	err = nft.outputHostBaseRules(c, nft.wanIface)
	if err != nil {
		return fmt.Errorf(`nft.outputHostBaseRules(%q): %w`, nft.wanIface, err)
	}
	err = nft.inputTrustIPSetRules(c, nft.wanIface)
	if err != nil {
		return fmt.Errorf(`nft.inputTrustIPSetRules(%q): %w`, nft.wanIface, err)
	}
	err = nft.outputTrustIPSetRules(c, nft.wanIface)
	if err != nil {
		return fmt.Errorf(`nft.outputTrustIPSetRules(%q): %w`, nft.wanIface, err)
	}
	err = nft.inputPublicRules(c, nft.wanIface)
	if err != nil {
		return fmt.Errorf(`nft.inputPublicRules(%q): %w`, nft.wanIface, err)
	}
	err = nft.outputPublicRules(c, nft.wanIface)
	if err != nil {
		err = fmt.Errorf(`nft.outputPublicRules(%q): %w`, nft.wanIface, err)
	}
	return err
}
