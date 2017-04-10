package api

import (
	"fmt"
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
func (c *APIClient) GetAllowSudoNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_99_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraAllowSudoScriptBody)
}

// GetAddIPCustomizeNoteID get add ip customize note id
func (c *APIClient) GetAddIPCustomizeNoteID(serverID string, ip string, subnet string) (string, error) {
	noteName := fmt.Sprintf("_30_%s_%d__", serverID, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraAddIPForEth1ScriptBodyFormat, ip, subnet)
	return c.getCustomizeNoteID(noteName, noteBody)
}

// GetChangeGatewayCustomizeNoteID get change gateway address customize note id
func (c *APIClient) GetChangeGatewayCustomizeNoteID(serverID string, gateway string) (string, error) {
	noteName := fmt.Sprintf("_20_%s_%d__", serverID, time.Now().UnixNano())
	noteBody := fmt.Sprintf(sakuraChangeDefaultGatewayScriptBody, gateway)
	return c.getCustomizeNoteID(noteName, noteBody)
}

// GetDisableEth0CustomizeNoteID get disable eth0 customize note id
func (c *APIClient) GetDisableEth0CustomizeNoteID(serverID string) (string, error) {
	noteName := fmt.Sprintf("_10_%s_%d__", serverID, time.Now().UnixNano())
	return c.getCustomizeNoteID(noteName, sakuraDisableEth0ScriptBody)
}

func (c *APIClient) getCustomizeNoteID(noteName string, noteBody string) (string, error) {

	existsNotes, err := c.client.Note.Reset().WithNameLike(noteName).Limit(1).Find()

	//すでに登録されている場合
	if len(existsNotes.Notes) > 0 {
		return fmt.Sprintf("%d", existsNotes.Notes[0].ID), nil
	}

	//ない場合はここで作成する
	newNote := c.client.Note.New()
	newNote.Name = noteName
	newNote.Content = noteBody

	note, err := c.client.Note.Create(newNote)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", note.ID), nil

}

// DeleteNote delete note
func (c *APIClient) DeleteNote(strId string) error {
	id, res := ToSakuraID(strId)
	if !res {
		return fmt.Errorf("NoteID is invalid: %v", strId)
	}

	_, err := c.client.Note.Delete(id)
	return err
}
