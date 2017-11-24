package com.letstalk.data_layer

import java.util.UUID

import com.letstalk.data_models.{ Message, UserModel }

import scala.collection.mutable

class MemoryChatStorage extends ChatStorage {

  private val messageBuffer: mutable.HashMap[UUID, Message] = mutable.HashMap()
  override def storeMessage(message: Message): Unit = {
    messageBuffer put (message.id, message)
  }

  override def retrieveMessage(id: UUID): Option[Message] = {
    messageBuffer get id
  }

  private val userBuffer: mutable.HashMap[UUID, UserModel] = mutable.HashMap()
  override def storeUser(user: UserModel): Unit = {
    userBuffer put (user.id, user)
  }

  override def retrieveUser(id: UUID): Option[UserModel] = {
    userBuffer get id
  }

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
