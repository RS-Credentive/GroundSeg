<script>
  import Clipboard from 'clipboard'
  import { onMount } from 'svelte'
  import { structure, registerServiceAgain } from '$lib/stores/websocket';
  import { devShipClass } from '$lib/stores/devclass'
  export let patp

  $: ship = ($structure?.urbits?.[patp]?.info) || {}
  $: memUsage = (ship?.memUsage) || 0
  $: diskUsage = (ship?.diskUsage) || 0
  $: loom = (ship?.loomSize) || 0
  $: loomActual = 2 ** loom / (1024 * 1024)
  $: svcRegStatus = (ship?.serviceRegistrationStatus) || "ok"
  $: vere = (ship?.vere) || ""

  $: chars = (patp.replace(/-/g,"").length) || 0
  $: shipClass = (chars == 3 ? "GALAXY" : chars == 6 ? "STAR" : chars == 12 ? "PLANET" : chars > 12 ? "MOON" : "UNKNOWN") || "ERROR"

  let copy
  let copied = false

  onMount(()=>{
    copy = new Clipboard('#patp');
    copy.on("success", ()=> {
      copied = true;
      setTimeout(()=> copied = false, 1000)
    })
  })

</script>

<div class="header">
  <div class="patp-wrapper">
    <div class="ship-class">{shipClass}
      <sup>{vere.toUpperCase()}</sup>
    </div>
    <div class="patp" id="patp" data-clipboard-text={patp}>
      {#if copied}
        COPIED!
      {:else}
        <!-- dev modification -->
        {#if $devShipClass == "star"}
          {patp.slice(0,6).toUpperCase()}
        {:else if $devShipClass == "planet"}
          {patp.slice(0,13).toUpperCase()}
        {:else}
          {patp.toUpperCase()}
        {/if}
      {/if}
    </div>
  </div>
  <div class="settings-wrapper">
    <div class="settings">
      <div class="settings-text">RAM</div>
      <div class="settings-val">{parseInt(memUsage/(1024*1024))} MB / {loomActual} MB</div>
    </div>
    <div class="settings">
      <div class="settings-text">DISK</div>
      <div class="settings-val">{(diskUsage/(1024*1024)).toFixed(2)} MB</div>
    </div>
  </div>
</div>

<style>
  .header {
    background-color: var(--bg-card);
    color: var(--text-card-color);
    position: absolute;
    height: 150px;
    width: calc(1173px - 150px);
    left: 150px;
    border-radius: 16px 16px 0 0;
    position: relative;
  }
  .patp-wrapper {
    background: var(--fg-card);
    position: absolute;
    top: 16px;
    left: 16px;
    width: 567px;
    height: 134px;
    border-radius: 16px;
  }
  .ship-class {
    font-family: var(--title-font);
    margin: 28px 0 12px 20px;
    font-size: 18px;
  }
  .patp {
    cursor: pointer;
    font-family: var(--title-font);
    margin-left: 20px;
    font-size: 32px;
  }
  .settings-wrapper {
    background: var(--fg-card);
    font-family: var(--title-font);
    position: absolute;
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 320px;
    padding: 24px 32px;
    right: 0;
    border-radius: 0px 16px;
  }
  .settings {
    display: flex;
    font-size: 24px;
  }
  .settings-text {
    flex: 1;
  }
  .btn {
    width: 30%;
    font-family: var(--regular-font);
    font-size: 12px;
    line-height: 32px;
    background-color: var(--btn-secondary);
    color: var(--text-card-color);
    border-radius: 8px;
    cursor: pointer;
  }
  .btn:hover {
    background: var(--bg-card);
  }
  .btn:disabled {
    pointer-events: none;
    opacity: .6;
  }
</style>
