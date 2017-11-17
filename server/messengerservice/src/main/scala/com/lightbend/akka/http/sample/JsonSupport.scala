package com.lightbend.akka.http.sample

import akka.http.scaladsl.marshallers.sprayjson.SprayJsonSupport
import com.lightbend.akka.http.sample.UserRegistryActor.ActionPerformed
import com.lightbend.akka.http.sample.data_models._
import com.lightbend.akka.http.sample.routes.MessageData
import spray.json.DefaultJsonProtocol

/**
  * Unmarshall json into the data models
  */
trait JsonSupport extends SprayJsonSupport {
  // import the default encoders for primitive types (Int, String, Lists etc)
  import DefaultJsonProtocol._

  implicit val actionPerformedJsonFormat = jsonFormat1(ActionPerformed)

  // Types of User info
  implicit val personalInfoJsonFormat = jsonFormat1(PersonalInfo)
  implicit val contactInfoJsonFormat = jsonFormat2(ContactInfo)

  // Types of User models
  implicit val normalUserJsonFormat = jsonFormat3(NormalUser)
  implicit val administratorUserJsonFormat = jsonFormat3(AdministratorUser)

  // Message Payload
  implicit val incomingMessagePayloadJsonFormat = jsonFormat2(IncomingMessagePayload)
  implicit val outgoingMessagePayloadJsonFormat = jsonFormat3(OutgoingMessagePayload)

  implicit  val messageDataJsonFormat = jsonFormat3(MessageData)

}
