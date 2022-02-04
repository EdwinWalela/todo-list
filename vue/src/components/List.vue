<template>
    <div>
        <div :class="[todo.isComplete ? 'complete':'']" id="todo-container" :key="todo.id" v-for="(todo,index) in todos" >
            
            <h3 >{{ todo.title }}</h3>
             <input @click="$emit('update-item',index)" id="check" type="checkbox"/>
             <button @click="$emit('delete-item',index)">x</button>
            <p>Due: {{ 
               todo.timestamp
                
            }} </p>
            
            
           
        </div>
    </div>
</template>

<script>
export default {
    name:'List',
    props:{
        todos:Array
    },
}
</script>

<style scoped>
    #todo-container{
        box-shadow: 0px 2px 3px rgba(0,0,0,0.3);
        padding: 5px 10px;
        margin:20px 0;
         border-left: solid 4px transparent;
    }

    h3{
        font-size: 0.9em;
        width: 81%;
        display: inline-block;
    }
    p{
        font-size: 0.8em;;
    }

    .complete {
       border-left: solid 4px rgb(0, 204, 0) !important;

    }

    button{
        display: inline-block;
        
        margin: 0em 1em;
        position: relative;
        bottom:5px
    }

    button:nth-of-type(1){
        border:none;
        background: red;
        color: white;
    }
     button:nth-of-type(2){
        border:none;
        background: green;
        color: white;
    }

    #check{
        padding:10px;
        height: 18px;
        width: 18px;
    }
    
</style>