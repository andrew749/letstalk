package com.letstalk.data_models

import slick.jdbc.PostgresProfile.api._

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * An immutable message object to encapsulate data being sent
 * @param sender
 * @param recipient
 * @param payload
 */
final case class Message(sender: UserModel, recipient: UserModel, payload: Option[MessagePayload])


/**
  * The table schema definition for a table.
  * @param tag
  */
class MessageTable(tag: Tag) extends Table[(Int, Int, Int, Int)](tag, "Messages") {
  def messageId = column[Int]("ID", O.PrimaryKey)
  def sender_id = column[Int]("SENDER_ID")
  def recipient_id = column[Int]("RECIPIENT_ID")
  def payload = column[Int]("PAYLOAD")

  // need to add this projection
  def * = (messageId, sender_id, recipient_id, payload)

  val payloads = TableQuery[MessagePayloadTable]
  def payloadForeign = foreignKey("PAYLOAD_FK", payload, payloads)(_.id)
}
