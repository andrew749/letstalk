package com.lightbend.akka.http.sample

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * Encapsulate personal information of a user.
 * @param name name of a user
 */
case class PersonalInfo(name: String)

/**
 * Encapsulate contact information for a user.
 * @param email email address of a user
 */
case class ContactInfo(email: String, phone: String)