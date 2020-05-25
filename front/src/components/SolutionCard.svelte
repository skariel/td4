<script>
    export let solution;

    function getRunTime(t) {
        if (!t.Valid) {
            return "not started yet";
        }
        else {
            return t.Time;
        }
    }
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

    .displayname {

    }

    .topbar {
        grid-column: 1 / 5;
        grid-row: 1;
        display: flex;
    }

    .solutionid {
        margin-left: auto;
    }

    .created {
        grid-column: 1 / 5;
        grid-row: 2;
    }

    .updated {
        grid-column: 1 / 5;
        grid-row: 3;
    }

    .status {
        grid-column: 1 / 5;
        grid-row: 4;
    }
    .started {
        grid-column: 1 / 5;
        grid-row: 5;
    }
    .ended {
        grid-column: 1 / 5;
        grid-row: 6;
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
    <h4 class="created">c: {solution.ts_created}</h4>
    <h4 class="updated">u: {solution.ts_updated}</h4>

    <h4 class="status">Status: {solution.status}</h4>

    <h4 class="started">s: {getRunTime(solution.ts_start)}</h4>
    <h4 class="ended">e: {getRunTime(solution.ts_end)}</h4>
</div>

