package com.letstalk.data_models

import slick.jdbc.PostgresProfile.api._

class MessageTable(tag: Tag) extends Table[(Int, Int, Int, Int)](tag, "Messages") {
  def messageId = column[Int]("ID", O.PrimaryKey)
  def sender = column[Int]("SENDER_ID")
  def recipient = column[Int]("RECIPIENT_ID")
  def payload = column[Int]("PAYLOAD")

  def * = (messageId, sender, recipient, payload)
}

