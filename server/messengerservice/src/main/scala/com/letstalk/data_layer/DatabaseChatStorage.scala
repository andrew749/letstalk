package com.letstalk.data_layer

import com.letstalk.data_models.{ Message, UserModel }

class DatabaseChatStorage extends ChatStorage {

  override def storeMessage(message: Message): Unit = {
    // FIXME: store in db
    ???
  }

  override def retrieveMessage(id: String): Option[Message] = {
    // FIXME: retrieve from db
    ???
  }

  override def storeUser(user: UserModel) = {
    // FIXME: store in db
    ???
  }

  override def retrieveUser(id: String): Option[UserModel] = {
    // FIXME: retrieve from db
    ???
  }

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
