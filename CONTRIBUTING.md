# Contributing to zk-Harness

### Table of Contents

- [Contributing to zk-Harness](#contributing-to-zk-harness)
    - [Table of Contents](#table-of-contents)
    - [General Guidelines](#general-guidelines)
    - [Reporting Bugs](#reporting-bugs)
      - [Before Submitting A Bug](#before-submitting-a-bug)
      - [How Do I Submit a (Good) Bug?](#how-do-i-submit-a-good-bug)
    - [Pull Requests](#pull-requests)
    - [How to integrate a new framework?](#how-to-integrate-a-new-framework)


### General Guidelines

Please do:

- **DO** include meaningful commit messages. You can find a guide to commit messages [here](https://github.com/RomuloOliveira/commit-messages-guide#good-practices).
- **DO** Before including new benchmarks for a circuit, make sure that the framework is already supported. If not, please file an issue to integrate the new framework, following the guidelines on [How to integrate a new framework?](#how-to-integrate-a-new-framework)
- **DO** open an issue for design discussion before making any major changes.
- **DO** read our [documentation](https://github.com/zkCollective/zk-Harness/tree/main/documentation) to understand how zk-Harness is designed.
- **DO** follow the conventions as described in the benchmarking integration of individual projects.
- **DO** give priority to the current style of the project or file you're
  changing even if it diverges from the general guidelines.
- **DO** include tests when adding new circuits.
- **DO** update README.md files in the source tree and other documents to be up
  to date with changes in the code.
- **DO** keep the discussions focused. When a new or related topic comes up it's
  often better to create a new issue than to side track the discussion.

### Reporting Bugs
#### Before Submitting A Bug 
* Ensure the bug is not already reported by searching on GitHub under 
[Issues](https://github.com/zkCollective/zk-Harness/issues).
#### How Do I Submit a (Good) Bug?
* If you are unable to find an open issue addressing the problem, open a new one. Be sure to include a 
**title and clear description**, as much relevant information as possible, and a **code sample** or 
an **executable test case** demonstrating the unexpected behavior.
* Describe the **exact steps** to **reproduce the problem** in as many details as possible. When 
listing steps, don't just say what you did, but explain how you did it. 
* Provide **specific examples** to demonstrate the steps. Include links to files or GitHub projects, or 
copy/pasteable snippets, which you use in those examples. If you're providing snippets in the issue, 
use [Markdown code blocks](https://help.github.com/articles/getting-started-with-writing-and-formatting-on-github/).
* Describe the **behavior you observed** after following the steps and explain the 
problem with that behavior.
* Explain the **behavior you expected** instead and why.
* **Can you reliably reproduce the issue?** If not, provide details about how often the problem 
happens and under which conditions it normally happens.

### Pull Requests

Pull requests will be reviewed by the project team against criteria including:
* purpose - is this change useful
* test coverage - are there unit/integration/acceptance tests demonstrating the change is effective
* code consistency - naming, comments, design
* changes that are solely formatting are likely to be rejected
* changes that do not follow the project design are likely to be rejected

Always write a clear log message for your commits. One-line messages are fine for small changes, but 
bigger changes should contain more detail.

### How to integrate a new framework?

The main principle zk-Harness aims for is allowing developers to benchmark custom circuit implementations without extensive overhead. Therefore, adding a new framework requires consideration and documentation.

To integrate a yet unsupported framework, please take the following steps as a guideline:

1. First, fork the ``zk-benchmarks`` repository
2. Create a ``./<framework_name>`` folder in the root folder of the repository.
3. Create a custom benchmarking script that *(i)* reads from the standardized input of the ``config.json`` as described in the ``config`` folder and outputs *(ii)* the standardized logs as described in the ``logs`` folder.
    * For example, benchmarking for ``gnark`` is done through a custom CLI, based on [cobra](https://github.com/spf13/cobra)
    * Your script should be able to take a variety of arguments as specified in the config.json, such that benchmarks can be easily executed and extended. E.g., a common command in the gnark integration would be ``./gnark groth16 --circuit=sha256 --input=_input/circuit/sha256/input_3.json --curve=bn254``
4. Modify the ``_scripts/reader/X`` scripts to include your newly created script as described in step 3, which is called if the ``project`` field of the respective config contains the ``<framework_name>`` of your newly added ZKP framework.
  * The ``_scripts/reader/X`` processing python scripts are invoked by ``__main__.py`` based on the ``category`` (currently arithmetics, ec, circuit) field as specified in the config.
5. Create a documentation in the ``./<framework_name>/tutorials`` folder to outline how others can include new circuits / benchmarks for another ``category`` in the framework. Depending on the project you pursue, it should contain documentation on either, or all, of the following:
    * How to add a new circuit implementation
    * How to run tests for integrated circuits
    * How to benchmark a new ``category`` in the ZKP-framework
6. Add config files for running the benchmarks in `_input/config/` and add a make rule for the new framework in the Makefile.

If you follow the specified interfaces for config and logs, your framework specific benchmarking should seamlessly integrate into the zk-Harness frontend.

Once finished, please create a Pull Request and assign it to one of the maintainers for review and correct implementation.