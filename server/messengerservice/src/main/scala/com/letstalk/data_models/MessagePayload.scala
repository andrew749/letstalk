package com.letstalk.data_models

import java.util.UUID

import slick.jdbc.PostgresProfile.api._

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

trait MessagePayload

/**
 * Encapsulate incoming message payload
 * @param messageContent The body of the message
 * @param time unix timestamp of the time this message was sent
 */
case class IncomingMessagePayload(messageContent: String, time: Long) extends MessagePayload {}

/**
 * Encapsulate outgoing message payload
 * @param id id of the message
 * @param messageContent payload
 * @param time that the message was acked by the server
 */
case class OutgoingMessagePayload(id: UUID, messageContent: String, time: Long) extends MessagePayload {}

class MessagePayloadTable(tag: Tag) extends Table[(UUID, String, Long)](tag, "MESSAGES") {
  def id = column[UUID]("ID", O.PrimaryKey)
  def messageContent = column[String]("DATA")
  def time = column[Long]("DATE_SENT")

  def * = (id, messageContent, time)
}
