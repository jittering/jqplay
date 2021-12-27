<script lang="ts">
  import { onMount } from "svelte";
  import { EditorState, EditorView, basicSetup } from "@codemirror/basic-setup";
  import { ViewUpdate } from "@codemirror/view";
  import { LanguageSupport } from "@codemirror/language";
  import { debounce } from "lodash-es";

  // props
  let classes: string = "";
  export let editorView: EditorView = null;
  export let langs: Array<LanguageSupport>;
  export let readOnly: boolean = false;
  export let options = {};
  export let value: string;
  export { classes as class };

  let lastValue: string = null;

  let element;

  onMount(() => createEditor(options));

  const onUpdate = EditorView.updateListener.of(
    debounce((v: ViewUpdate) => {
      if (!v.docChanged) {
        return;
      }
      const newVal = v.state.doc.toString();
      if (newVal !== value) {
        lastValue = newVal; // ensure we don't trigger a state update on editing
        value = newVal;
      }
    }, 250)
  );

  $: if (element) {
    createEditor(options);
  }

  function createEditor(options) {
    if (!window) return;
    if (!element) return;

    if (editorView) element.innerHTML = "";

    editorView = new EditorView({
      parent: element,
    });

    updateState(value, langs);
  }

  function updateState(value: string, langs: Array<LanguageSupport>) {
    if (lastValue === value) {
      return;
    }
    let editorState = EditorState.create({
      doc: value,
      extensions: [
        basicSetup,
        onUpdate,
        ...langs,
        EditorView.editable.of(!readOnly),
      ],
    });
    editorView.setState(editorState);
    lastValue = value;
  }

  $: if (editorView) {
    updateState(value, langs);
  }
</script>

<div bind:this={element} class={classes} />

<style global lang="postcss">
  /* BASICS */
  :global(.cm-editor) {
    height: var(--cm-height, 300px);
    direction: ltr;
    color: var(--cm-text-color);
    background: var(--cm-background-color);
  }

  :global(.cm-editor) {
    .cm-scroller {
      font-family: "SF Mono", "DejaVu Sans Mono", Menlo, Monaco, Consolas,
        Courier, monospace !important;
    }
  }
</style>
