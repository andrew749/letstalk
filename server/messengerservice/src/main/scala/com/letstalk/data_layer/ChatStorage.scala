package com.letstalk.data_layer

import com.letstalk.data_models.Message

trait ChatStorage {
  def storeMessage(message: Message)
  def retrieveMessage(id: String): Option[Message]
  def storeUser()
  def storeUserInfo()
  def storeContactInfo()
}
