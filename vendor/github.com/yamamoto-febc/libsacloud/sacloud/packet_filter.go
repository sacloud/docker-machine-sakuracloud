package sacloud

import (
	"fmt"
	"regexp"
)

// PacketFilter type of PacketFilter
type PacketFilter struct {
	*Resource
	Index       int `json:",omitempty"`
	Name        string
	Description string `json:",omitempty"`

	//HACK API呼び出しルートにより数字/文字列が混在する
	// PackerFilterのCREATE時は文字列、以外は数値となる。現状利用しないためコメントとしておく
	// RequiredHostVersion int    `json:",omitempty"`

	Notice     string                    `json:",omitempty"`
	Expression []*PacketFilterExpression `json:",omitempty"`
}

func AllowPacketFilterProtocol() []string {
	return []string{"tcp", "udp", "icmp", "fragment", "ip"}
}

type PacketFilterExpression struct {
	Protocol        string `json:",omitempty"`
	SourceNetwork   string `json:",omitempty"`
	SourcePort      string `json:",omitempty"`
	DestinationPort string `json:",omitempty"`
	Action          string `json:",omitempty"`
	Description     string `json:",omitempty"`
}

func CreateNewPacketFilter() *PacketFilter {
	return &PacketFilter{
		Expression: []*PacketFilterExpression{},
	}
}

func (p *PacketFilter) ClearRules() {
	p.Expression = []*PacketFilterExpression{}
}

func (p *PacketFilter) AddTCPRule(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool) error {

	err := p.validatePort(sourcePort)
	if err != nil {
		return err
	}
	err = p.validatePort(destPort)
	if err != nil {
		return err
	}

	exp := &PacketFilterExpression{
		Protocol:        "tcp",
		SourceNetwork:   sourceNetwork,
		SourcePort:      sourcePort,
		DestinationPort: destPort,
		Action:          p.getActionString(isAllow),
		Description:     description,
	}

	p.Expression = append(p.Expression, exp)
	return nil
}

func (p *PacketFilter) AddUDPRule(sourceNetwork string, sourcePort string, destPort string, description string, isAllow bool) error {

	err := p.validatePort(sourcePort)
	if err != nil {
		return err
	}
	err = p.validatePort(destPort)
	if err != nil {
		return err
	}

	exp := &PacketFilterExpression{
		Protocol:        "udp",
		SourceNetwork:   sourceNetwork,
		SourcePort:      sourcePort,
		DestinationPort: destPort,
		Action:          p.getActionString(isAllow),
		Description:     description,
	}

	p.Expression = append(p.Expression, exp)
	return nil
}

func (p *PacketFilter) AddICMPRule(sourceNetwork string, description string, isAllow bool) error {

	exp := &PacketFilterExpression{
		Protocol:      "icmp",
		SourceNetwork: sourceNetwork,
		Action:        p.getActionString(isAllow),
		Description:   description,
	}

	p.Expression = append(p.Expression, exp)
	return nil
}

func (p *PacketFilter) AddFragmentRule(sourceNetwork string, description string, isAllow bool) error {

	exp := &PacketFilterExpression{
		Protocol:      "fragment",
		SourceNetwork: sourceNetwork,
		Action:        p.getActionString(isAllow),
		Description:   description,
	}

	p.Expression = append(p.Expression, exp)
	return nil
}

func (p *PacketFilter) AddIPRule(sourceNetwork string, description string, isAllow bool) error {

	exp := &PacketFilterExpression{
		Protocol:      "ip",
		SourceNetwork: sourceNetwork,
		Action:        p.getActionString(isAllow),
		Description:   description,
	}

	p.Expression = append(p.Expression, exp)
	return nil
}

func (p PacketFilter) getActionString(isAllow bool) string {
	action := "deny"
	if isAllow {
		action = "allow"
	}
	return action
}

func (p *PacketFilter) validatePort(expression string) error {
	if expression == "" {
		return nil

	}

	match, err := regexp.MatchString("^[0-9]*$", expression)
	if err != nil {
		return err
	}
	if match {
		return nil
	}

	match, err = regexp.MatchString("^[0-9]{1,5}-[0-9]{1,5}$", expression)
	if err != nil {
		return err
	}
	if match {
		return nil
	}

	return fmt.Errorf("Bad syntax:%s", expression)
}
