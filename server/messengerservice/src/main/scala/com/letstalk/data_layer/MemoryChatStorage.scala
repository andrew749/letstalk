package com.letstalk.data_layer

import java.util.UUID
import scala.collection.mutable

import akka.actor.Actor

import com.letstalk.data_models.{ Message, Thread, UserModel }

class MemoryChatStorage extends ChatStorage {

  private val messageBuffer: mutable.HashMap[UUID, Message] = mutable.HashMap()
  override def storeMessage(message: Message): Unit = {
    println(s"storing message ${message.payload}")
    messageBuffer put (message.id, message)
  }

  override def retrieveMessage(id: UUID): Option[Message] = {
    messageBuffer get id
  }

  override def retrieveMessages(threadId: UUID): Seq[Message] = {
    messageBuffer.values.filter(_.threadId == threadId).toSeq
  }

  private val threadBuffer: mutable.HashMap[UUID, Thread] = mutable.HashMap()
  override def storeThread(thread: Thread): Unit = {
    threadBuffer put (thread.id, thread)
  }

  override def retrieveThread(id: UUID): Option[Thread] = {
    threadBuffer get id
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
