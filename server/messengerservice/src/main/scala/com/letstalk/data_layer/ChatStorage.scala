package com.letstalk.data_layer

trait ChatStorage  {
  def storeMessage()
  def storeMessagePayload()
  def storeUser()
  def storeUserInfo()
  def storeContactInfo()
}
