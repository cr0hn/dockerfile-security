Dockerfile-sec
==============

Usage
-----

With remote rules
+++++++++++++++++

.. code-block:: console

    $ dockerfile-sec -r http://127.0.0.1:9999/rules/credentials.yaml Dockerfile

With local rules
++++++++++++++++

.. code-block:: console

    $ dockerfile-sec -r rules/credentials.yaml Dockerfile

With many rules files
+++++++++++++++++++++

.. code-block:: console

    $ dockerfile-sec -r rules/java.yaml -r rules/credentials.yaml Dockerfile

Export results as json
++++++++++++++++++++++

.. code-block:: console

    $ dockerfile-sec -o results.json -r rules/java.yaml -r rules/credentials.yaml Dockerfile

Quiet mode
++++++++++

Don't output nothing into console

.. code-block:: console

    $ dockerfile-sec -q Dockerfile

Filtering false positives
+++++++++++++++++++++++++

**By ignore file**

Dockerfile-sec allows to ignore rules by using a file that contains the rules you want to ignore.

.. code-block:: console

    $ dockerfile-sec -F ignore-rules Dockerfile

Ignore file format contains the *IDs* of rules you want to ignore. **one ID per line**. Example:

.. code-block:: console

    core-001
    core-007

**By cli**

You also can use cli to ignore specific *IDs*:

.. code-block:: console

    $ dockerfile-sec -i core-001,core007 Dockerfile

Using as pipeline
+++++++++++++++++

You also can use dockerfile-sec as UNIX pipeline.

Loading Dockerfile from stdin:

.. code-block:: console

    $ cat Dockerfile | dockerfile-sec -i core-001,core007

Exposing results via pipe:


.. code-block:: console

    $ cat Dockerfile | dockerfile-sec -i core-001,core007 | jq



Output formats
--------------

JSON Output format
++++++++++++++++++

.. code-block:: json

    [
      {
        "description": "Missing content trust",
        "id": "core-001",
        "reference": "https://snyk.io/blog/10-docker-image-security-best-practices/",
        "severity": "Low"
      },
      {
        "description": "Missing USER sentence in dockerfile",
        "id": "core-002",
        "reference": "https://snyk.io/blog/10-docker-image-security-best-practices/",
        "severity": "Medium"
      }
    ]


References
----------

- https://snyk.io/blog/10-docker-image-security-best-practices/
- https://medium.com/microscaling-systems/dockerfile-security-tuneup-166f1cdafea1
- https://medium.com/@tariq.m.islam/container-deployments-a-lesson-in-deterministic-ops-a4a467b14a03