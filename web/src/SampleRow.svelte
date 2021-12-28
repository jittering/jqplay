<script lang="ts">
  import { Button } from "smelte";
  import { createEventDispatcher } from "svelte";

  export let sample: any;

  const dispatch = createEventDispatcher();
</script>

<tr class="border-t">
  <td>
    <code>{sample.code}</code>
  </td>
  <td>{sample.text}</td>
  <td>
    <Button
      icon="assignment"
      text
      light
      flat
      on:click={() => dispatch("sample", { sample })}
    />
  </td>
</tr>
