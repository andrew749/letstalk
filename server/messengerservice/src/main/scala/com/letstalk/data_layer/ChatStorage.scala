package com.letstalk.data_layer

import com.letstalk.data_models.{ Message, UserModel }

trait ChatStorage {
  def storeMessage(message: Message)
  def retrieveMessage(id: String): Option[Message]
  def storeUser(user: UserModel)
  def retrieveUser(id: String): Option[UserModel]
  def storeUserInfo()
  def storeContactInfo()
}
