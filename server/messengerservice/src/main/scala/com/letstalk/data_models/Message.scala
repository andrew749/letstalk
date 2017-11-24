package com.letstalk.data_models

import java.util.UUID

import slick.jdbc.PostgresProfile.api._

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * An immutable message object to encapsulate data being sent
 * @param id: unique identifier for the message
 * @param sender: User object corresponding to who is sending the message.
 * @param recipient: User object corresponding to who is receiving the message.
 * @param payload: Message data
 */
final case class Message(id: UUID, sender: UserModel, recipient: UserModel, payload: Option[MessagePayload])

/**
 * The table schema definition for a table.
 * @param tag
 */
class MessageTable(tag: Tag) extends Table[(UUID, UUID, UUID, UUID)](tag, "Messages") {
  def messageId = column[UUID]("ID", O.PrimaryKey)
  def senderId = column[UUID]("SENDER_ID")
  def recipientId = column[UUID]("RECIPIENT_ID")
  def payload = column[UUID]("PAYLOAD")

  // need to add this projection
  def * = (messageId, senderId, recipientId, payload)

  val payloads = TableQuery[MessagePayloadTable]
  def payloadForeign = foreignKey("PAYLOAD_FK", payload, payloads)(_.id)
}
