package com.letstalk

/**
 * Unmarshall json into the data models
 */
trait JsonSupport extends SprayJsonSupport with DefaultJsonProtocol {
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
