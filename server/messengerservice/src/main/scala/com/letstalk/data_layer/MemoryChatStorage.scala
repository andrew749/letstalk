package com.letstalk.data_layer

import com.letstalk.data_models.Message

import scala.collection.mutable

class MemoryChatStorage extends ChatStorage {

  private val messageBuffer: mutable.HashMap[String, Message] = mutable.HashMap()
  override def storeMessage(message: Message): Unit = {
    messageBuffer.put(message.id, message)
  }

  override def retrieveMessage(id: String): Option[Message] = {
    messageBuffer get id
  }

  override def storeUser(): Unit = ???

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
