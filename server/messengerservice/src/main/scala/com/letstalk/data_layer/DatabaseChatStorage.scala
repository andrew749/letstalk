package com.letstalk.data_layer

import scala.collection.mutable

trait DatabaseChatStorage extends ChatStorage { this: DataManager =>

  implicit var dataLayers: mutable.Buffer[ChatStorage]

  dataLayers += this


  override def storeMessage(): Unit = ???

  override def storeMessagePayload(): Unit = ???

  override def storeUser(): Unit = ???

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
