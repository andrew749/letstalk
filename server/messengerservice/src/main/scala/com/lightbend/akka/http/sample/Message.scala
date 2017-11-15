package com.lightbend.akka.http.sample

/**
 * Created by andrewcodispoti on 2017-11-14.
 */
class Message(id: String, payload: MessagePayload) {

  def getPayload() = {
    payload
  }

}

case class IncomingMessage(id: String, sender: UserModel, payload: MessagePayload) extends Message(id, payload = payload) {}
case class OutgoingMessage(id: String, recipient: UserModel, payload: MessagePayload) extends Message(id, payload = payload) {}
