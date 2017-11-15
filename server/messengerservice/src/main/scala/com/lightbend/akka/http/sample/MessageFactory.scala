package com.lightbend.akka.http.sample
import akka.actor.{ Actor, ActorLogging, Props }

/**
 * Created by andrewcodispoti on 2017-11-14.
 */
object MessageFactory {
  final case class ParseJson(jsonString: String)

  def props: Props = Props[MessageFactory]
}

class MessageFactory extends Actor with ActorLogging {
  override def receive: Receive = {
    case _ => new IncomingMessage(
      id = "1",
      payload = new MessagePayload("Hello World"),
      sender = NormalUser(
        info = PersonalInfo("Andrew"),
        contactInfo = ContactInfo("andrewcod749@gmail.com", "5555555555")
      )
    )
  }
}
