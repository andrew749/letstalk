package com.letstalk.data_layer

import java.util.UUID

import akka.actor.{ Actor, ActorLogging, Props }
import com.letstalk.data_models.{ Message, Thread, UserModel }

import scala.collection.mutable

trait MessageEvent

case class GetMessages(threadId: UUID) extends MessageEvent
case class GetThread(id: UUID) extends MessageEvent

case class Messages(values: Seq[Message])

/**
 * This class receives messages with data that should be stored in
 * the appropriate datalayer type for it.
 *
 * Massive WIP
 */

class DataManager(storage: ChatStorage) extends Actor with ActorLogging {

  def receive: Receive = {
    case message: Message =>
      storage storeMessage message

    case GetMessages(id) =>
      println(s"get messages with id ${id}")
      val messages = storage retrieveMessages id
      sender() ! Messages(messages)

    case thread: Thread =>
      storage storeThread thread

    case GetThread(id) =>
      val thread = storage retrieveThread id
      sender() ! thread
  }

}

object DataManager {
  def props(storage: ChatStorage): Props = {
    Props(classOf[DataManager], storage)
  }
}
