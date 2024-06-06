Analyze the k6 test script and extract the extensions that the script depends on.

**Sources**

Dependencies can come from three sources: k6 test script, manifest file, `K6_DEPENDENCIES` environment variable.

Primarily, the k6 test script is the source of dependencies. The test script and the local and remote JavaScript modules it uses are recursively analyzed. The extensions used by the test script are collected. In addition to the require function and import expression, the `"use k6 ..."` directive can be used to specify additional extension dependencies. If necessary, the `"use k6 ..."` directive can also be used to specify version constraints.

       "use k6>0.49";
       "use k6 with k6/x/faker>=0.2.0";
       "use k6 with k6/x/toml>v0.1.0";
       "use k6 with xk6-dashboard*";

Dependencies and version constraints can also be specified in the so-called manifest file. The default name of the manifest file is `package.json` and it is automatically searched from the directory containing the test script up to the root directory. The `dependencies` property of the manifest file contains the dependencies in JSON format.

    {"dependencies":{"k6":">0.49","k6/x/faker":">=0.2.0","k6/x/toml":>v0.1.0","xk6-dashboard":"*"}}

Dependencies and version constraints can also be specified in the `K6_DEPENDENCIES` environment variable. The value of the variable is a list of dependencies in a one-line text format.

       k6>0.49;k6/x/faker>=0.2.0;k6/x/toml>v0.1.0;xk6-dashboard*

**Format**

By default, dependencies are written as a JSON object. The property name is the name of the dependency and the property value is the version constraints of the dependency.

    {"k6":">0.49","k6/x/faker":">=0.2.0","k6/x/toml":>v0.1.0","xk6-dashboard":"*"}

Additional output formats:

 * `text` - One line text format. A semicolon-separated sequence of the text format of each dependency. The first element of the series is `k6` (if there is one), the following elements follow each other in lexically increasing order based on the name.

       k6>0.49;k6/x/faker>=0.2.0;k6/x/toml>v0.1.0;xk6-dashboard*

 * `js` - A consecutive, one-line JavaScript string directives. The first element of the series is `k6` (if there is one), the following elements follow each other in lexically increasing order based on the name.

       "use k6>0.49";
       "use k6 with k6/x/faker>=0.2.0";
       "use k6 with k6/x/toml>v0.1.0";
       "use k6 with xk6-dashboard*";

**Output**

By default, dependencies are written to standard output. By using the `-o/--output` flag, the dependencies can be written to a file.
