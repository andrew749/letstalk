package com.letstalk.data_layer

import java.util.UUID

import com.letstalk.data_models.{ Message, Thread, UserModel }

trait ChatStorage {
  def storeMessage(message: Message)
  def retrieveMessage(id: UUID): Option[Message]
  def retrieveMessages(threadId: UUID): Seq[Message]
  def storeThread(thread: Thread)
  def retrieveThread(id: UUID): Option[Thread]
  def storeUser(user: UserModel)
  def retrieveUser(id: UUID): Option[UserModel]
  def storeUserInfo()
  def storeContactInfo()
}
