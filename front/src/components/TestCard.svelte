<script>
  import { timeSince } from "../routes/utils";

  export let test;

  const ts = timeSince(new Date(test.ts_updated));
</script>

<style>
  .main {
    display: grid;
    background-color: rgb(236, 230, 226);
    border: 3px solid #ffffff;
    border-radius: 7px;
    padding-top: 10px;
    padding-left: 10px;
    padding-right: 10px;
    padding-bottom: 0px;
    grid-template-columns: 25% 25% 25% 25%;
    grid-template-rows: auto auto auto auto auto;
  }

  .main:hover {
    border: 3px solid rgb(255, 62, 0);
    padding-top: 10px;
    padding-left: 10px;
    padding-right: 10px;
    padding-bottom: 0px;
  }

  .avatar {
    width: 40px;
    height: 40px;
    margin-right: 20px;
    margin-top: 3px;
  }

  .topbar {
    grid-column: 1 / 5;
    grid-row: 1;
    display: flex;
    border-bottom: black solid 1px;
  }

  .testid {
    margin-left: auto;
    font-size: 150%;
    margin-top: 7px;
  }

  .testtitle {
    grid-column: 1 / 5;
    grid-row: 2;
    word-wrap: break-word;
    margin-top: 7px;
  }

  .testdescr {
    grid-column: 1 / 5;
    grid-row: 3;
    word-wrap: break-word;
    color: #777777;
  }

  .updated {
    grid-column: 1 / 5;
    grid-row: 4;
    font-size: 80%;
  }

  .failed {
    grid-column: 1;
    grid-row: 5;
    display: flex;
    flex-direction: column;
    text-align: center;
    margin-top: 7px;
    font-size: 80%;
  }

  .passed {
    grid-column: 2;
    grid-row: 5;
    display: flex;
    flex-direction: column;
    text-align: center;
    margin-top: 7px;
    font-size: 80%;
  }

  .pending {
    grid-column: 3;
    grid-row: 5;
    display: flex;
    flex-direction: column;
    text-align: center;
    margin-top: 7px;
    font-size: 80%;
  }

  .wip {
    grid-column: 4;
    grid-row: 5;
    display: flex;
    flex-direction: column;
    text-align: center;
    margin-top: 7px;
    font-size: 80%;
  }

  .bottom {
    grid-column: 1/5;
    grid-row: 5;
    border-top: black solid 1px;
  }
</style>

<!-- <div class="main" onclick="location.href='{"/test?id="+test.id+"&page=0"}';"> -->
<div class="main">

  <div class="topbar">
    <img class="avatar" src={test.avatar} alt="avatar" />
    <div style="display:flex; flex-direction: column; width: 100%">
      <h4 class="displayname">{test.display_name}</h4>
      <h4 class="updated">{ts}</h4>
    </div>
    <h4 class="testid">#{test.id}</h4>
  </div>

  {#if test.descr.title < 128}
    <h4 class="testtitle">
        <a href={'/test?id=' + test.id + '&page=0'}>{test.title}</a>
    </h4>
  {:else}
    <h4 class="testtitle">
        <a href={'/test?id=' + test.id + '&page=0'}>{test.title.slice(0,128)+'...'}</a>
    </h4>
  {/if}

  {#if test.descr.length < 256}
    <h4 class="testdescr">{test.descr}</h4>
  {:else}
    <h4 class="testdescr">{test.descr.slice(0,256)+'...'} <a href={'/test?id=' + test.id + '&page=0'}>more</a> </h4>
  {/if}

  <div class="failed">
    <h4>{test.total_fail}</h4>
    <h4>fail</h4>
  </div>
  <div class="passed">
    <h4>{test.total_pass}</h4>
    <h4>pass</h4>
  </div>
  <div class="pending">
    <h4>{test.total_pending}</h4>
    <h4>todo</h4>
  </div>
  <div class="wip">
    <h4>{test.total_wip}</h4>
    <h4>wip</h4>
  </div>
  <div class="bottom" />

</div>
