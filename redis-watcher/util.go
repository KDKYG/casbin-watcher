package rediswatcher

import (
	"log"
)

type CallbackFunc func(msg string, update, updateForAddPolicy, updateForRemovePolicy, updateForRemoveFilteredPolicy, updateForSavePolicy func(string, interface{}))

func CustomDefaultFunc(defaultFunc func(string, interface{})) CallbackFunc {
	return func(msg string, update, updateForAddPolicy, updateForRemovePolicy, updateForRemoveFilteredPolicy, updateForSavePolicy func(string, interface{})) {
		msgStruct := &MSG{}
		err := msgStruct.UnmarshalBinary([]byte(msg))
		if err != nil {
			log.Println(err)
		}
		invoke := func(f func(string, interface{})) {
			if f == nil {
				f = defaultFunc
			}
			f(msgStruct.ID, msgStruct.Params)
		}
		switch msgStruct.Method {
		case "Update":
			invoke(update)
		case "UpdateForAddPolicy":
			invoke(updateForAddPolicy)
		case "UpdateForRemovePolicy":
			invoke(updateForRemovePolicy)
		case "UpdateForRemoveFilteredPolicy":
			invoke(updateForRemoveFilteredPolicy)
		case "UpdateForSavePolicy":
			invoke(updateForSavePolicy)
		}
	}
}
