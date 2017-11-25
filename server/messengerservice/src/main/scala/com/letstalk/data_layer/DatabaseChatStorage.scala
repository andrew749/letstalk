package com.letstalk.data_layer

import java.util.UUID

import com.letstalk.data_models.{ Message, Thread, UserModel }

class DatabaseChatStorage extends ChatStorage {

  override def storeMessage(message: Message): Unit = {
    // FIXME: store in db
    ???
  }

  override def retrieveMessage(id: UUID): Option[Message] = {
    // FIXME: retrieve from db
    ???
  }

  override def retrieveMessages(threadId: UUID): Seq[Message] = {
    // FIXME: retrieve from db
    ???
  }

  override def storeThread(thread: Thread) = {
    // FIXME: store in db
    ???
  }

  override def retrieveThread(id: UUID): Option[Thread] = {
    // FIXME: retrieve from db
    ???
  }

  override def storeUser(user: UserModel) = {
    // FIXME: store in db
    ???
  }

  override def retrieveUser(id: UUID): Option[UserModel] = {
    // FIXME: retrieve from db
    ???
  }

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
