package api

import (
	"encoding/json"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
	"time"
)

const (
	sakuraAllowSudoScriptBody = `#!/bin/bash

  # @sacloud-once
  # @sacloud-desc ubuntuユーザーがsudo出来るように/etc/sudoersを編集します
  # @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
  # @sacloud-require-archive distro-debian
  # @sacloud-require-archive distro-ubuntu

  export DEBIAN_FRONTEND=noninteractive
	echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers || exit 1
	sh -c 'sleep 10; shutdown -h now' &
  exit 0`

	sakuraAddIPForEth1ScriptBodyFormat = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: setup ip address for eth1
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	echo "auto eth1" >> /etc/network/interfaces
	echo "iface eth1 inet static" >> /etc/network/interfaces
	echo "address %s" >> /etc/network/interfaces
	echo "netmask %s" >> /etc/network/interfaces
	ifdown eth1; ifup eth1
	exit 0`

	sakuraChangeDefaultGatewayScriptBody = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: change default gateway
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	sed -i 's/gateway/#gateway/g' /etc/network/interfaces
	echo "up route add default gw %s" >> /etc/network/interfaces
	exit 0`

	sakuraDisableEth0ScriptBody = `#!/bin/bash

	# @sacloud-once
	# @sacloud-desc docker-machine-sakuracloud: disable eth0
	# @sacloud-desc （このスクリプトは、DebianもしくはUbuntuでのみ動作します）
	# @sacloud-require-archive distro-debian
	# @sacloud-require-archive distro-ubuntu

	export DEBIAN_FRONTEND=noninteractive
	sed -i 's/iface eth0 inet static/iface eth0 inet manual/g' /etc/network/interfaces
	ifdown eth0 || exit 0
	exit 0`
)

// GetAllowSudoNoteID get ubuntu customize note id
// FIXME
// workaround for [Non root ssh create sudo can't get password](https://github.com/docker/machine/issues/1569)
// [PR #1586](https://github.com/docker/machine/pull/1586)がマージされるまで暫定
// スクリプト(Note)を使ってubuntuユーザがsudo可能にする
func (c *Client) GetAllowSudoNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_99_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraAllowSudoScriptBody)
}

// GetAddIPCustomizeNoteID get add ip customize note id
func (c *Client) GetAddIPCustomizeNoteID(serverID string, ip string, subnet string) (string, error) {
	noteName := fmt.Sprintf("_30_%s_%d__", serverID, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraAddIPForEth1ScriptBodyFormat, ip, subnet)
	return c.getCustomizeNoteID(noteName, noteBody)
}

// GetChangeGatewayCustomizeNoteID get change gateway address customize note id
func (c *Client) GetChangeGatewayCustomizeNoteID(serverID string, gateway string) (string, error) {
	noteName := fmt.Sprintf("_20_%s_%d__", serverID, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraChangeDefaultGatewayScriptBody, gateway)
	return c.getCustomizeNoteID(noteName, noteBody)
}

// GetDisableEth0CustomizeNoteID get disable eth0 customize note id
func (c *Client) GetDisableEth0CustomizeNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_10_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraDisableEth0ScriptBody)
}

func (c *Client) getCustomizeNoteID(noteName string, noteBody string) (string, error) {

	var (
		method = "GET"
		uri    = "note"
		body   = sakura.Request{
			Filter: map[string]interface{}{"Name": noteName},
		}
	)
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	uri = fmt.Sprintf("%s?%s", uri, bodyJSON)
	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var existsNote sakura.SearchResponse
	if err := json.Unmarshal(data, &existsNote); err != nil {
		return "", err
	}

	//すでに登録されている場合
	if existsNote.Count > 0 {
		return existsNote.Notes[0].ID, nil
	}

	//ない場合はここで作成する
	method = "POST"
	uri = "note"
	n := sakura.Request{
		Note: &sakura.Note{
			Name:    noteName,
			Content: noteBody,
		},
	}

	data, err = c.newRequest(method, uri, n)
	if err != nil {
		return "", err
	}
	var s sakura.Response
	if err := json.Unmarshal(data, &s); err != nil {
		return "", err
	}

	return s.Note.ID, nil

}

// DeleteNote delete note
func (c *Client) DeleteNote(id string) error {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("note/%s", id)
	)

	_, err := c.newRequest(method, uri, nil)
	if err != nil {
		return err
	}
	return nil
}
