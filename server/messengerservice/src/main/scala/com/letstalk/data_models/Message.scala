package com.letstalk.data_models

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
final case class Message(id: String, sender: UserModel, recipient: UserModel, payload: Option[MessagePayload])

/**
 * The table schema definition for a table.
 * @param tag
 */
class MessageTable(tag: Tag) extends Table[(String, String, String, String)](tag, "Messages") {
  def messageId = column[String]("ID", O.PrimaryKey)
  def sender_id = column[String]("SENDER_ID")
  def recipient_id = column[String]("RECIPIENT_ID")
  def payload = column[String]("PAYLOAD")

  // need to add this projection
  def * = (messageId, sender_id, recipient_id, payload)

  val payloads = TableQuery[MessagePayloadTable]
  def payloadForeign = foreignKey("PAYLOAD_FK", payload, payloads)(_.id)
}
