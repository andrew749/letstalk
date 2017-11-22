package com.letstalk.data_layer

import com.letstalk.data_models.Message

class DatabaseChatStorage extends ChatStorage {

  override def storeMessage(message: Message): Unit = {
    // FIXME: store in db
  }

  override def retrieveMessage(id: String): Option[Message] = {
    // FIXME: retrieve from db
    ???
  }

  override def storeUser(): Unit = ???

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
