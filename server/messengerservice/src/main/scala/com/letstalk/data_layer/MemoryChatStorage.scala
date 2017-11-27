package com.letstalk.data_layer

import java.util.UUID
import scala.collection.mutable

import akka.actor.Actor

import com.letstalk.data_models.{ Message, Thread, UserModel }

class MemoryChatStorage extends ChatStorage {

  private val messageBuffer: mutable.HashMap[UUID, Message] = mutable.HashMap()
  override def storeMessage(message: Message): Unit = {
    messageBuffer put (message.id, message)
  }

  override def retrieveMessage(id: UUID): Option[Message] = {
    messageBuffer get id
  }

  override def retrieveMessages(threadId: UUID): Seq[Message] = {
    messageBuffer.values.filter(_.threadId == threadId).toSeq
  }

  private var threadUserBuffer = mutable.HashMap[UUID, List[UUID]]()
  private var userThreadBuffer = mutable.HashMap[UUID, List[UUID]]().withDefaultValue(Nil)
  override def storeThread(thread: Thread): Unit = {
    threadUserBuffer put (thread.id, thread.userIds)
    thread.userIds.foreach((userId: UUID) =>
      userThreadBuffer.update(userId, thread.id :: userThreadBuffer(userId)))
  }

  override def retrieveThreads(userId: UUID): Seq[Thread] = {
    println(userThreadBuffer)
    userThreadBuffer.get(userId).getOrElse(Seq()).map((threadId: UUID) =>
      Thread(threadId, threadUserBuffer.get(threadId).getOrElse(Nil)))
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
