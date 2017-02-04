# Pomenator

This is a tiny utility I wrote for myself to deal with hellish struggle
of publishing artifacts to Maven Central. As far as I can tell, Maven
Central is the defacto standard location that Open Source stuff needs to
be published to in order for it to be usable for the community. And as
far as I can tell, Nexus / Sonatype seems to be the only possibility to
access Maven Central.

Apart from venting my anger, this rant aims to provide an accessible
introduction to publishing a project to Maven Central. It also
documents a tool I wrote to automate the process.

I'm sure the fine folks at Sonatype are dedicated people and are
aware of the shortcomings of their system, so my apologies to them for
the following rant. I am an ungrateful little whiny bitch that should
appreciate the fact they're hosting all this infrastructure.

The system seems to be intended as a cruel joke, it's documented
through a series of unlinked youtube
[videos](http://central.sonatype.org/pages/producers.html) that force
you to listen to 2 minutes of awful flamenco music at the start of each
of the ~10 clips and are narrated by the slowest speaking moderator ever
who also has a very thick accent. Alternatively, there is a long
introduction on how to configure maven. No words lost on the conceptual
background of the whole process or what actually happens below the hood.

The process is as follows:

- create a Jira account (https://issues.sonatype.org)
- create an issue to register a new groupid (this is the groupid from
  the maven groupId - artifactId - version trifecta), this issue needs
  to be created in the Project: "Community Support - Open Source Project
  Repository Hosting"
- wait until someone approves the issue. This seems to be automated if 
  your email domain corresponds with the group id

Next, you need to:

- Log in to https://oss.sonatype.org a web frontend was obviously
  built using some sort of web gui builder that SAP discontinued in the
  early 90's.
- Don't try to figure out what any of the word crap there means. Seriously.
- Click on the link "Staging Upload" (in the left navigation)

![Navigation](https://raw.githubusercontent.com/a2800276/pomenator/master/doc/sc_nav.png)

- Select the "Upload mode" you're using (more on this later, for now,
  let's assume you have a pom and a bunch of signed jars.) Select "POM"
- Now upload the POM file (this is done using semi standard html widgets
  which for whatever reason assigns your POM a new fake name rooted at
  'c:\')
- Now you need to upload 7 (!) further files:
  - a jar file containing your classes
  - a jar file containing your sources
  - a jar file containing the javadoc generated from the sources
  - for each of the jar files and the POM file you need to upload
    a PGP signature (*.ASC) (see below)

![Upload](https://raw.githubusercontent.com/a2800276/pomenator/master/doc/sc_upload.png)

- You do this by selecting each of the 7 files individually and then
  clicking "Add Artifact", finally hit: "Upload Artifacts"

This is probably a good time to take up smoking. Because -- you've
guessed it -- you're not done yet.

- Click on "Staging Repositories" (in the left navigation)
- Don't get distracted by the fact that there's a million other people's
  repositories (that's just what they call them, don't try to make sense
  of the name "repository") in the list.

![Other People's Repositories](https://raw.githubusercontent.com/a2800276/pomenator/master/doc/sc_opp.png)

- Your repository, i.e. your upload is way at the bottom of the list.
  Scroll down. Keep scrolling.
- All the millions of other Repos are named `central_bundle-12345`,
  yours will be called `{your_group_id}-1234`
- If you don't find your upload on the list, it's likely that the
  mainframe is down or has not processed the batch yet. You'll need to
  click on "Refresh"
- Keep clicking on "Refresh"
- Once it shows up, you can click on your project and a bunch of fancy
  html widget trees and shit appear. You can click around on them and
  you may find an indication of an error, in case you've missed
  something. It's very likely you make a mistake the first couple of
  times you try to release.

![Super Fancy Tree](https://raw.githubusercontent.com/a2800276/pomenator/master/doc/sc_super_fancy_tree.png)

- If there are no errors in the fancy html tree widgets, the "Release"
  button at the top of the page should no longer be greyed out. You can
  click on that.
- Then you need to "refresh" to make sure your upload has disappeared
  from the list.

This is probably a good time to learn assembler or take up heroin.
Because you're only almost done. You still need to go back to Jira
(https://issues.sonatype.org) and activate synchronization. *You do this
my adding a comment to the original issue you used to ask for a
group_id*. 

The comment part is a one time only thing. Once the synchronization is
set up, your uploads will eventually (this is dependant on when the guys
at sonatype have a chance to run down to the Staples to buy new
punchcards for their mainframe) show up on maven central. Sooner or
later, your stuff should show up for searches on search.maven.org and
the actual artefacts will be located at

    http://repo1.maven.org/maven2/{group_seperated_by_slashes}/{artifact_id}/{version}/

## OMGWTFBBQ you can't be serious!?

I most certainly am. Kind of. Some possibilities exist to "simplify"
this process:

### Use maven

Maven has tasks or whatever maven calls tasks to handle all of this
crap. If you're using maven already, great. No way in hell I'm
completely restructuring my project to use maven to manage my project.

### gradle or sbt or use clojure

Ok. whatever

### Use bundle uploads

In case you are stuck doing this manually, here's the single best piece
of advice I can offer you: sonatype/nexus accepts "bundle" uploads. This
is one of the alternative modes that's offered in the "Upload Mode"
chooser. This allows you to pack all 8 of the jars and signatures you
need to upload into a single jar file. Rumor has it that Sonatype
implemented this option because the cost of buying new mice was becoming
unsustainable for them.

### Automate this yourself using the REST API provided by sonatype

This is what I'm trying to do here. Unfortunately -- I shit you not --
Sonatype does not provide any documentation for their REST API (unless
you have an Admin account) and their official (well blog) advice is
this: read [this blog
post](http://blog.sonatype.com/2012/07/learning-the-nexus-rest-api-read-the-docs-or-fire-up-a-browser/)
which whines about how hard it is to write up-to-date documentation
for a REST API and then explains how to use the Chrome network inspector
to see what REST-calls their web frontend makes. I wish I were joking.

Here is all the documentation you need:

- create the POM, jars and signatures as outlined above.
- pack them all together into a single jar
- upload them to: https://oss.sonatype.org/service/local/staging/bundle_upload
- The upload is a POST with a raw body, don't need none of that mime-multipart crud
- All HTTP requests can be authenticated with HTTP Basic Authentication
- A successful upload returns status 201 and a bunch of JSON that
  contains an id in a property name `stagedRepositoryIds`
- Finally, you "release" your upload by POSTing to: 
  https://oss.sonatype.org/service/local/staging/bulk/promote
- which contains a Content-Type header indicating Javascript
- and contains a body (again, raw) like this:
    
      {"data":{"autoDropAfterRelease":true,"description":"{WHATEVER}","stagedRepositoryIds":["{YOUR_REPO_ID}"]}}

- If you try to call the release endpoint too early, you'll get some
  sort of 500 error indicating that the repo is still "open" or
  "transitioning". I assume this means that the folks at sonatype are at
  the blacksmith's getting iron rings made for their core memory.
- It seems that waiting 10 seconds between the upload and the
  'release' call is usually sufficient to resolve the 500 problem.

### Use this tool

Wouldn't suggest that. It's a crappy tool I wrote for myself and your
computer will probably catch fire. And you'll just be adding yet another
crappy tool to your (no doubt already horrifically convoluted) release
toolchain.

The following notes are therefore for myself:

- create a .json file corresponding to the required pom data. A
  [sample](https://github.com/a2800276/pomenator/blob/master/test/bootstrap.json)
  is in the `test` directory.
- additionally, the json file contains the following properties:
  - sources : an array of directories containing sources
  - classes : an array of directories containing classes
  - output  : the output and working directory where all the artifacts
    get dumped.

The main tool is the `bin/pomenator` tool. It is written in Go. I
[started out](https://github.com/a2800276/pomerator) writing it in Java. Java is fine, but it wasn't worth the
effort to learn how to do an HTTP POST using the built in
HttpUrlConnection. Running the tool gives you a set of flags you can use:

    Usage of ./pomenator:
      -config string
        	name of the main config
      -key-id string
        	pgp key id
      -keyring string
        	secure keyring file
      -pgp-passwd string
        	secure keyring password
      -repo-passwd string
        	sonatype passwd
      -repo-user string
        	sonatype username
      -secrets string
        	json file containing pgp keyring, id and password, and your 
          sonatype id, passwd (default "./.secrets.json")

The `-config` flag is mandatory and refers to the config described
above.  The other flags allow you to enter your super secret passwords
and stuff. You can store them all in a json file that looks like this:

    {
       "secring" : "./secring.gpg"
      ,"keyId" : "A4B924E5"
      ,"passwd" : "<super secret pgp keyring password>"
      ,"repo_user" : "<sonatype user>"
      ,"repo_passwd" : "<sonatype password>"
    }

You can also mix and match passing secrets via flags and via file.
Should you do this, values provided by flags take precedence over those
in the config file.

The example above assumes you're using a local keyring (see below on how
to set this up). But you can also use a system wide keyring (probably
located in `~/.gnupg/secring.gpg`) You'll need to replace the key id
with your own key id, the value "A4B924E5" above is an example.

Note that Nexus allows you to use a special replacement "token" username
and password so that your real username and password aren't involved.
It's probably a good idea to do this, because:

- the tokenised uid/pswd don't allow access to Jira
- you can't reset your password in Nexus m) ...
- if you're at all like me, you'll end up publishing your
  `.secrets.json` to github sooner or later

To get a token, click on username->Profile at the top right corner of nexus
(oss.sonatype.org).

![User Profile](Profil://raw.githubusercontent.com/a2800276/pomenator/master/doc/sc_token_profile.png )

Then choose "User Token" in the drop down that appears and
click the "Access User Token" button. Use this instead of your real
username and password.

![User Token](Profil://raw.githubusercontent.com/a2800276/pomenator/master/doc/sc_token_token.png )

Also note that while the main pom config allows you to define a number
of different artifacts/projects in a single file, all of them will be
signed with the same pgp key and uploaded and released using the same
sonatype credentials.

## GPG

Now all that's left is dealing with all the GPG crap. Thanksfully it's a
one time setup thing. Here's the short version:

- install GPG, you will need to figure out how to do this yourself
- create an openpgp signing key for your project. My advice is to use
  the following invocation:

      gpg --no-default-keyring \
          --keyring ./pubring.gpg \
          --secret-keyring ./secring.gpg \
          --gen-key

- Note your `gpg` command may be called `gpg2`.
- `--no-default-keyring` means to not store the new key in the system
  wide keyring with all of your other gpg keys. 
- `--keyring` tells gpg where to store the public half of the key
  instead
- just as `--secret-keyring` tells gpg where to store the secret half.
  Make sure that file is added to your `.gitignore` or similar so you
  don't publish it. On the other hand, who cares? it's not like anyone
  checks signatures.
- The filenames provided must have at least one slash ('/') in them,
  else gpg will store them in some magic secret hidden directory off
  your home directory and you'll be left wondering why generating the
  key didn't work.

You'll be asked a million questions during key gneration. Most of them
are obvious. Here's my advice:  

when asked about the 'kind of key' that I want:

    Please select what kind of key you want:
      (1) RSA and RSA (default)
      (2) DSA and Elgamal
      (3) DSA (sign only)
      (4) RSA (sign only)

I pick one of the "sign only" options to make sure these keys aren't
ever used to encrypt email.

When asked for the keysize, I pick the maximum available.

When asked for the validity period, I pick the default.

Whether or not you store the key on the default keyring or in a
dedicated file is up to you. If you have a bunch of projects all signed
with the same key, it may be easier not to copy keyrings around.

Now you need to figure out the key id:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        --list-keys

Your keyring should only contain a single key, so you get output like:

    ./pubring.gpg
    ----------------
    pub   3072D/4B8FD19D 2017-02-03
    uid       [ultimate] DONT TRUST THIS KEY (TEST SIGNING KEY DONT TRUST SOFTWARE SIGNED WITH THIS KEY!) <test@example.com>

This keys key id is `4B8FD19D`. You still need to publish the key:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        --send-keys 4B8FD19D

Of course you need to substitute your own key id. You can also provide a
custom keyserver to publish the key on using the `--keyserver` argument
as sonatype suggests. I have no idea why they suggest this, or if they
enforce it. Anyway, just to keep you on your toes:
`hkp://pool.sks-keyservers.net` is the key server they'd like you to use.

Finally, to sign the jars generating the ASC files that you need for
Sonatype to accept your uploads:


    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --secret-keyring ./secring.gpg \
        -u 4B8FD19D \ 
        -ab whatever_your.jar

Obviously taking care to replace `whatever_your.jar` with the name of
the jar file you want to sign and using your own key id. This will
generate a `whatever_your.jar.asc` file.

If you'd like you can check that the generated signature is correct:

    gpg --no-default-keyring \
        --keyring ./pubring.gpg \
        --verify whatever_your.jar.asc

You'll get output indicating whether the signature corresponds to the
file. It will say "BAD" in capital letters if it doesn't.

That said, the tool doesn't generate and publish keys, but it does take
care of all the signature generation once you're set up, so the signing
stage above is only in case you want to do hings manually and nobody
ever verifies signatures anyway so just forget about the `--verify`
option.




