package com.lightbend.akka.http.sample.data_models

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * Base class of a message data.
 *
 * @param id The id of this specific message
 * @param messageContent The body of the message
 */
case class MessagePayload(id: String, messageContent: String, time: java.util.Date) {}
