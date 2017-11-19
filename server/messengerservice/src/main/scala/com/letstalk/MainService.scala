package com.letstalk

//#main-class
object MainService extends App
    with UserRoutes
    with MessageRoutes {

  // set up ActorSystem and other dependencies here
  //#main-class
  val config = ConfigFactory.load()

  //#server-bootstrapping
  implicit val system: ActorSystem = ActorSystem("LetsTalk", config)
  implicit val materializer: ActorMaterializer = ActorMaterializer()
  //#server-bootstrapping

  private val log = Logging.getLogger(system, this)

  // Needed for the Future and its methods flatMap/onComplete in the end
  implicit val executionContext: ExecutionContext = system.dispatcher

  val timeout = Timeout(5 seconds)

  val userRegistryActor: ActorRef = system.actorOf(UserRegistryActor.props, "userRegistryActor")

  /**
    * FIXME: Remove this. just here to create fake data
    */
  def createUsers() = {
    val pinfo = PersonalInfo("Andrew")
    val pinfo2 = PersonalInfo("Andrew 2")
    val cinfo = ContactInfo("test@gmail.com", "5555555555")
    userRegistryActor ! CreateUser(NormalUser("acod", pinfo, cinfo))
    userRegistryActor ! CreateUser(NormalUser("andrew", pinfo2, cinfo))
    log.debug("Created fake users")
  }

  // create fake users
  createUsers()


  //#main-class
  // from the UserRoutes trait
  lazy val routes: Route = concat(
    userRoutes,
    messageRoute
  )


  //#main-class

  //#http-server
  val serverBindingFuture: Future[ServerBinding] = Http().bindAndHandle(routes, "localhost", 8080)

  log.info(s"Server online at http://localhost:8080/\nPress RETURN to stop...")

  StdIn.readLine()

  serverBindingFuture
    .flatMap(_.unbind())
    .onComplete { done =>
      done.failed.map { ex => log.error(ex, "Failed unbinding") }
      system.terminate()
    }
  //#http-server
  //#main-class
}
