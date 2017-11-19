package com.letstalk.sample.data_models

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
case class OutgoingMessagePayload(id: String, messageContent: String, time: Long) extends MessagePayload {}
