<script>
	import { timeSince } from '../routes/utils';

    export let solution;

    const ts = timeSince(new Date(solution.ts_updated));

</script>

<style>
    .main {
        display:       grid;
        background-color: rgb(236,230,226);
        border:        3px solid #ffffff;
        border-radius: 7px;
        padding:       10px;
        grid-template-columns: 25% 25% 25% 25%;
        grid-template-rows: auto auto auto auto auto;
    }

    .main:hover {
        border:        3px solid rgb(255,62,0);
        padding: 10px;
    }

    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;
    }

    .topbar {
        grid-column: 1 / 5;
        grid-row: 1;
        display: flex;
    }

    .solutionid {
        margin-left: auto;
    }

    .updated {
        grid-column: 1 / 5;
        grid-row: 3;
    }

    .status {
        grid-column: 1 / 5;
        grid-row: 4;
    }
</style>

<!-- <div class="main" onclick="location.href='{"/test?id="+test.id+"&page=0"}';"> -->
<div class="main">

    <div class="topbar">
        <img class="avatar" src={solution.avatar} alt="avatar"/>
        <h4 class="displayname">{solution.display_name}</h4>
        <h4 class="solutionid">
            <a href={"/solution?id="+solution.id}>
                solution {solution.id}
            </a>
        </h4>
    </div>
    <h4 class="updated">updated: {ts}</h4>
    <h4 class="status">Status: {solution.status}</h4>
</div>

