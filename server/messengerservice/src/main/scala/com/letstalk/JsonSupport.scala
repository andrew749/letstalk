package com.letstalk

import java.util.UUID

import akka.http.scaladsl.marshallers.sprayjson.SprayJsonSupport
import com.letstalk.UserRegistryActor.ActionPerformed
import com.letstalk.data_models._
import com.letstalk.routes.{ MessageData, SendMessageResponse, CreateThread }
import spray.json._

import scala.collection.immutable

/**
 * Unmarshall json into the data models
 */
trait JsonSupport extends SprayJsonSupport with DefaultJsonProtocol {

  /**
   * Custom serializer for UUID
   * @return
   */
  implicit def UUIDFormat = new RootJsonFormat[UUID] {
    override def write(obj: UUID) = {
      new JsString(obj.toString)
    }

    override def read(json: JsValue) = {
      UUID.fromString(json.convertTo[String])
    }
  }
  // import the default encoders for primitive types (Int, String, Lists etc)
  implicit val actionPerformedJsonFormat: RootJsonFormat[ActionPerformed] = jsonFormat1(ActionPerformed)

  // Types of User info
  implicit val personalInfoJsonFormat: RootJsonFormat[PersonalInfo] = jsonFormat1(PersonalInfo)
  implicit val contactInfoJsonFormat: RootJsonFormat[ContactInfo] = jsonFormat2(ContactInfo)

  // Types of User models
  implicit val normalUserJsonFormat: RootJsonFormat[NormalUser] = jsonFormat3(NormalUser)
  implicit val administratorUserJsonFormat: RootJsonFormat[AdministratorUser] = jsonFormat3(AdministratorUser)

  // Message Payload
  implicit val incomingMessagePayloadJsonFormat = jsonFormat2(IncomingMessagePayload)
  implicit val outgoingMessagePayloadJsonFormat = jsonFormat3(OutgoingMessagePayload)
  implicit object MessagePayloadFormat extends RootJsonFormat[MessagePayload] {
    override def write(obj: MessagePayload) = obj match {
      case msg: IncomingMessagePayload => msg.toJson
      case msg: OutgoingMessagePayload => msg.toJson
    }
    override def read(value: JsValue) = value match {
      case obj: JsObject if (obj.fields.size == 2) => value.convertTo[IncomingMessagePayload]
      case obj: JsObject => value.convertTo[OutgoingMessagePayload]
    }
  }

  implicit val messageDataJsonFormat: RootJsonFormat[MessageData] = jsonFormat3(MessageData)
  implicit val sendMessageResponseFormat: RootJsonFormat[SendMessageResponse] =
    jsonFormat1(SendMessageResponse)

  implicit val messageJsonFormat: RootJsonFormat[Message] = jsonFormat4(Message)
  implicit val threadJsonFormat: RootJsonFormat[Thread] = jsonFormat2(Thread)
  implicit val createThreadJsonFormat: RootJsonFormat[CreateThread] = jsonFormat1(CreateThread)
}
