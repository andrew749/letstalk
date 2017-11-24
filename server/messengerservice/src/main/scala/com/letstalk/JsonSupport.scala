package com.letstalk

import java.util.UUID

import akka.http.scaladsl.marshallers.sprayjson.SprayJsonSupport
import com.letstalk.UserRegistryActor.ActionPerformed
import com.letstalk.data_models._
import com.letstalk.routes.MessageData
import spray.json.{ DefaultJsonProtocol, JsString, JsValue, JsonFormat, RootJsonFormat }

/**
 * Unmarshall json into the data models
 */
trait JsonSupport extends SprayJsonSupport with DefaultJsonProtocol {

  /**
    * Custom serializer for UUID
    * @return
    */
  implicit def UUIDFormat = new JsonFormat[UUID] {
    override def write(obj: UUID) = {
      new JsString(obj.toString)
    }

    override def read(json: JsValue) = {
      UUID.fromString(json.toString)
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
  implicit val incomingMessagePayloadJsonFormat: RootJsonFormat[IncomingMessagePayload] = jsonFormat2(IncomingMessagePayload)
  implicit val outgoingMessagePayloadJsonFormat: RootJsonFormat[OutgoingMessagePayload] = jsonFormat3(OutgoingMessagePayload)

  implicit val messageDataJsonFormat: RootJsonFormat[MessageData] = jsonFormat3(MessageData)

}
