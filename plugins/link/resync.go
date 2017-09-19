package link

import (
	"fmt"
	"time"
)

//unableToSync adds a link to the outOfSync list for the given peer
func (man *LinkManager) unableToSync(peerId peerID, linkId linkID) {
	if peerId == man.peerId {
		panic("unableToSync() called with self as peer.  Are you sure you want to do that?")
	}

	man.resyncMutex.Lock()
	defer man.resyncMutex.Unlock()

	if len(man.outOfSync) == 0 {
		go man.resyncProcess()
	}

	linkIdMap, have := man.outOfSync[peerId]
	if !have {
		linkIdMap = make(map[linkID]bool)
		man.outOfSync[peerId] = linkIdMap
	}
	linkIdMap[linkId] = true
}

func (man *LinkManager) resyncProcess() {
	isOutOfSync := true
	for isOutOfSync {
		time.Sleep(time.Second)
		fmt.Println("Attempting to resync...")

		man.resyncMutex.Lock()
		for peerId, linkIdMap := range man.outOfSync {
			linkIds := make([]linkID, len(linkIdMap))
			i := 0
			for linkId := range linkIdMap {
				linkIds[i] = linkId
				i++
			}

			if err := man.tryResyncUnsafe(peerId, linkIds); err != nil {
				fmt.Println(peerId, err)
			} else {
				delete(man.outOfSync, peerId)
				fmt.Println(peerId, "OK")
			}
		}

		if len(man.outOfSync) == 0 {
			isOutOfSync = false
		}
		man.resyncMutex.Unlock()
	}
	fmt.Println("All nodes back in sync.")
}