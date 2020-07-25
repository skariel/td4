<script>
	import { timeSince } from '../routes/utils';

    export let solution;

    const ts = timeSince(new Date(solution.ts_updated));

    let back_color = '';
    if (solution.status=='fail') {
        back_color = 'failed_back';
    }
    else if (solution.status=='pass') {
        back_color = 'passed_back';
    }


</script>

<style>
    .main {
        display:       grid;
        background-color: rgb(236,230,226);
        border:        3px solid #ffffff;
        border-radius: 7px;
        padding-top: 10px;
        padding-left: 10px;
        padding-right: 10px;
        padding-bottom: 0px;
        grid-template-columns: 25% 25% 25% 25%;
        grid-template-rows: auto auto auto auto auto;
    }

    .failed_back {
        background-color: rgb(130,55,55);
    }

    .passed_back {
        background-color: rgb(55,130,55);
    }

    .main:hover {
        border:        3px solid rgb(255,62,0);
        padding-top: 10px;
        padding-left: 10px;
        padding-right: 10px;
        padding-bottom: 0px;
    }

    .avatar {
        width: 40px;
        height: 40px;
        margin-right: 20px;
        margin-top: 0px;
    }

    .topbar {
        grid-column: 1 / 5;
        grid-row: 1;
        display: flex;
        border-bottom: black solid 1px;
    }

    .solutionid {
        margin-left: auto;
        font-size: 150%;
        margin-top: 7px;
    }

    .updated {
        grid-column: 1 / 5;
        grid-row: 3;
        font-size: 80%;
    }

    .displayname {
        margin-bottom: 1px;
    }

    .status {
        grid-column: 1 / 3;
        grid-row: 4;
        margin-top: 7px;
    }
    .status_light {
        grid-column: 3 / 5;
        grid-row: 4;
        margin-top: 7px;
        margin-left:auto;
        border-radius: 17px;
        margin-bottom: 7px;
        width: 21px;
    }

</style>

<!-- <div class="main" onclick="location.href='{"/test?id="+test.id+"&page=0"}';"> -->
<div class="main">

  <div class="topbar">
    <img class="avatar" src={solution.avatar} alt="avatar" />
    <div style="display:flex; flex-direction: column; width: 100%">
      <h4 class="displayname">{solution.display_name}</h4>
      <h4 class="updated">{ts}</h4>
    </div>
    <h4 class="solutionid"><a href={"/solution?id="+solution.id}>#{solution.id}</a></h4>
  </div>

   <h4 class="status">Status: {solution.status}</h4>
   <div class="status_light {back_color}" />
</div>

