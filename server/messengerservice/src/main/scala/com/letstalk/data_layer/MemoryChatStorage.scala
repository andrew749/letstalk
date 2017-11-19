package com.letstalk.data_layer

import com.letstalk.data_models.UserModel

import scala.collection.mutable
import scala.collection.mutable.ArrayBuffer

trait MemoryChatStorage extends ChatStorage { this: DataManager =>
  private var log = ArrayBuffer[UserModel]()

  implicit var dataLayers: mutable.Buffer[ChatStorage]

  override def storeMessage(): Unit = {}

  override def storeMessagePayload(): Unit = ???

  override def storeUser(): Unit = ???

  override def storeUserInfo(): Unit = ???

  override def storeContactInfo(): Unit = ???
}
