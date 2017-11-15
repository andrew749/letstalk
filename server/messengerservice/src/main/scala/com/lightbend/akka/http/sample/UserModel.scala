package com.lightbend.akka.http.sample

/**
 * Created by andrewcodispoti on 2017-11-14.
 */

/**
 * Base class to represent a user
 * @param info personal information to identify a user
 * @param contactInfo
 */
class UserModel(info: PersonalInfo, contactInfo: ContactInfo) {}

/**
 * Specific user of the platform
 * @param info personal information to identify a user
 * @param contactInfo
 */
case class NormalUser(info: PersonalInfo, contactInfo: ContactInfo) extends UserModel(info, contactInfo) {}

/**
 * Administrator of the platform
 *
 * Need a differentiation since an administrator account can perform special actions and have access to special interface.
 * @param info personal information to identify a user
 * @param contactInfo
 */
case class AdministratorUser(info: PersonalInfo, contactInfo: ContactInfo) extends UserModel(info, contactInfo) {}