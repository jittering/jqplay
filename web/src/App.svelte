<script lang="ts">
  import "smelte/src/tailwind.css";

  import { TextField, Checkbox, Chip, ProgressLinear, Button } from "smelte";
  import Service, { Jq } from "./service";
  import { debounce } from "lodash-es";

  import { json as langJson } from "@codemirror/lang-json";
  import CodeMirror from "./CodeMirror.svelte";
  import { onMount } from "svelte";
  import Panel from "./Panel.svelte";
  import Cheatsheet from "./Cheatsheet.svelte";
  import { samplesLeft, samplesRight } from "./samples";
  import { samplesJmesLeft, samplesJmesRight } from "./samples_jmespath";

  export let global: any;
  const jq = new Jq();
  let result = "";
  let jqVersion = "...";
  let commandLine = "";

  // switches for jq/jmespath
  let mode = "jq"; // jq or JMESPath
  let switchLabel = "JMESPath";
  let samples = { Left: samplesLeft, Right: samplesRight };

  const langs = [langJson()];

  const slurp = jq.getOpt("slurp");
  const nullInput = jq.getOpt("null-input");
  const compactOutput = jq.getOpt("compact-output");
  const rawInput = jq.getOpt("raw-input");
  const rawOutput = jq.getOpt("raw-output");

  const service = new Service();
  service.getJqVersion().then((ver) => {
    jqVersion = ver;
  });

  onMount(() => {
    loadInitialJSON();
  });

  function loadInitialJSON() {
    // load initial json, if avail
    service.getJqInput().then((json) => {
      if (!json) {
        json = "";
      }
      jq.j = json;
    });
  }

  function onClickDocumentation() {
    if (mode === "jq") {
      window.open("https://stedolan.github.io/jq/manual/", "_blank");
    } else {
      window.open("https://jmespath.org/tutorial.html", "_blank");
    }
  }

  function onClickSwitch() {
    // toggle the mode
    if (mode === "jq") {
      mode = "JMESPath";
      samples = { Left: samplesJmesLeft, Right: samplesJmesRight };
      switchLabel = "jq";
    } else {
      mode = "jq";
      samples = { Left: samplesLeft, Right: samplesRight };
      switchLabel = "JMESPath";
    }
  }

  function prettyPrint(input: string): string {
    const d = JSON.parse(input);
    return JSON.stringify(d, null, 2);
  }

  function onClickPrettyPrint() {
    try {
      jq.j = prettyPrint(jq.j);
    } catch (e) {
      // TODO: show err
    }
  }

  // run jq on input/filter/option changes
  $: {
    onChangeJq(jq);
  }

  const onChangeJq = debounce((jq) => {
    startProgressBar();
    if (mode === "jq") {
      return service
        .runJq(jq)
        .then((output) => {
          progress = -1;
          clearTimeout(progressTimeoutId);
          result = output;
        })
        .then(() => {
          service.getJqCommandLine(jq).then((output) => {
            commandLine = output;
          });
        });
    }

    // jmes
    return service.runJmesPath(jq).then((output) => {
      progress = -1;
      clearTimeout(progressTimeoutId);
      if (output) {
        if (output.startsWith("failed")) {
          result = output;
        } else {
          result = prettyPrint(output);
        }
      } else {
        result = "error: JMESPath search retruned null response";
      }
    });
  }, 250);

  // progressbar
  let progress = -1;
  let progressTimeoutId = null;
  $: {
    if (progress < 0 && progressTimeoutId !== null) {
      clearTimeout(progressTimeoutId);
    }
  }
  function startProgressBar() {
    progress = 0;
    function next() {
      progressTimeoutId = setTimeout(() => {
        if (progress < 0) {
          return;
        }
        if (progress === 100) {
          progress = 0;
        }
        progress += 1;
        next();
      }, 25);
    }
    next();
  }

  // load cheatsheet sample
  function loadSample(sample) {
    jq.j = sample.input_j;
    jq.q = sample.input_q;
  }
</script>

<main class="flex flex-col h-full">
  <div class="nav flex-initial container max-w-none p-2 flex">
    <div class="inline-block align-middle flex-initial">
      <a href="/" class="navbar-brand"
        ><img src="images/logo.png" alt="jqplay" /></a
      >
    </div>
    <div class="inline-block text-sm align-middle ml-4 flex-initial">
      A playground for
      {#if mode === "jq"}<a
          href="https://stedolan.github.io/jq/"
          class="navbar-link">jq</a
        >
        {jqVersion}
      {:else}
        <a href="https://jmespath.org/" class="navbar-link">JMESPath</a>
      {/if}
    </div>
    <div class="flex-auto" />
    <p class="inline-block docs flex-initial">
      <Chip icon="article" outlined on:click={onClickSwitch}
        >Switch to {switchLabel}</Chip
      >
      <Chip icon="article" outlined on:click={onClickDocumentation}
        >documentation</Chip
      >
    </p>
  </div>

  {#if progress >= 0}
    <ProgressLinear {progress} />
  {/if}

  <div class="flex-auto grid grid-cols-2 gap-4 overflow-hidden h-full p-4">
    <div class="inputs flex flex-col overflow-hidden h-full">
      <div class="flex-initial">
        {#if mode === "jq"}
          <h6>JQ Filter</h6>
        {:else}
          <h6>JMESPath Expression</h6>
        {/if}
        <TextField bind:value={jq.q} outlined />
      </div>
      <div class="flex-initial flex">
        <h6>JSON</h6>
        <div class="flex-auto" />
        <Button
          icon="format_align_left"
          small
          text
          flat
          title="Pretty Print"
          on:click={onClickPrettyPrint}
        />
        <Button
          icon="refresh"
          small
          text
          flat
          title="Reload initial JSON (if any)"
          on:click={loadInitialJSON}
        />
      </div>
      <CodeMirror
        class="json_input flex-auto h-full overflow-auto pt-2"
        bind:value={jq.j}
        {langs}
      />
    </div>

    <div class="outputs flex flex-col overflow-hidden h-full">
      <h6 class="">Result</h6>
      {#if mode === "jq"}
        <div class="jq_options flex-initial flex">
          <Checkbox
            label="Compact Output"
            bind:checked={compactOutput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Null Input"
            bind:checked={nullInput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Raw Input"
            bind:checked={rawInput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Raw Output"
            bind:checked={rawOutput.enabled}
            on:change={onChangeJq(jq)}
          />
          <Checkbox
            class="ml-2"
            label="Slurp"
            bind:checked={slurp.enabled}
            on:change={onChangeJq(jq)}
          />
        </div>
      {/if}
      <CodeMirror
        class="json_output flex-auto h-full overflow-auto"
        bind:value={result}
        {langs}
        readOnly={true}
      />
    </div>
  </div>

  {#if mode === "jq"}
    <Panel label="Command Line">
      <div class="text-center">
        <code>{commandLine}</code>
      </div>
    </Panel>
  {/if}
  <Panel label="Cheatsheet" collapsible={true}>
    <Cheatsheet {samples} on:sample={(ev) => loadSample(ev.detail.sample)} />
  </Panel>
</main>

<style lang="postcss">
  .nav {
    background-color: #222;
    color: #9d9d9d;
  }

  h6 {
    font-weight: bolder;
  }

  /* reduce checkbox padding */
  :global {
    .docs button {
      background-color: transparent;
    }
    .jq_options {
      label {
        padding-left: 0 !important;
      }
    }
  }

  code {
    color: #c7254e;
    background-color: #f9f2f4;
    border-radius: 4px;
  }

  :global {
    .cheatsheet {
      th:first-child,
      td:first-child {
        border-right-width: 0;
      }
    }
  }
</style>
