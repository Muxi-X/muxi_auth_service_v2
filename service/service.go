package service

import (
	"github.com/Muxi-X/muxi_auth_service_v2/model"
)

func CheckUserExisted(username, email string) bool {
	// 声明用于检查邮箱、用户名是否重复的通信信道；用于标识检查流程是否结束的信道
	sameEmailChannel, sameUsernameChannel, done := make(chan bool), make(chan bool), make(chan struct{})
	defer close(sameEmailChannel)
	defer close(sameUsernameChannel)
	// 自动关闭done信道，这样就可以以return来代替close()方法
	defer close(done)

	// 并发检查邮箱
	go func(email string) {
		_, err := model.GetUserByEmail(email)
		// 判断检查是否已经结束
		select {
		case <-done:
			return
		default:
			{
				// 没有结束，将结果输入信道
				if err != nil { // email not found
					sameEmailChannel <- true
				} else {
					sameEmailChannel <- false
				}
			}
		}
	}(email)

	// 并发检查同户名
	go func(username string) {
		_, err := model.GetUserByUsername(username)
		// 判断检查是否已经结束
		select {
		case <-done:
			return
		default:
			{
				// 没有结束，将结果输入信道
				if err != nil { // user not found
					sameUsernameChannel <- true
				} else {
					sameUsernameChannel <- false
				}
			}
		}
	}(username)

	// 用于标识用户名和邮箱重复检查的状态，false为没有重复
	userExisted := false

	// 最多循环两次
	for round := 0; !userExisted && round < 2; round++ {
		select {
		case emailResult := <-sameEmailChannel:
			{
				if !emailResult {
					userExisted = true
					// 不再等待另一个信道
					break
				}
			}
		case usernameResult := <-sameUsernameChannel:
			{
				if !usernameResult {
					userExisted = true
					// 不再等待另一个信道
					break
				}
			}
		}
	}
	return userExisted
}

func CheckUserNotExisted(username string) *model.UserModel {
	var user *model.UserModel
	checkChannel, done := make(chan *model.UserModel, 2), make(chan struct{})
	defer close(checkChannel)
	defer close(done)

	go func(username string) {
		checkUsernameUser, checkUsernameError := model.GetUserByUsername(username)
		select {
		case <-done:
			return
		default:
			if checkUsernameError != nil {
				checkChannel <- nil
			} else {
				checkChannel <- checkUsernameUser
			}
		}
	}(username)

	go func(email string) {
		checkEmailUser, checkEmailError := model.GetUserByEmail(email)
		select {
		case <-done:
			return
		default:
			if checkEmailError != nil {
				checkChannel <- nil
			} else {
				checkChannel <- checkEmailUser
			}
		}
	}(username)

	for user = range checkChannel {
		if user != nil {
			break
		}
	}

	return user
}
