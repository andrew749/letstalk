package com.letstalk.data_layer

import com.letstalk.sample.data_models.{ Message, UserModel }

import scala.collection.mutable.ArrayBuffer

class MemoryChatStorage extends ChatStorage {
  private var log = ArrayBuffer[UserModel]()
  def receive: Receive = {
    case msg @ Message(sender, recipient, payload) =>
      // add this message to our memory transcipt
      // FIXME: Add actual storage
      log += sender
  }

}
