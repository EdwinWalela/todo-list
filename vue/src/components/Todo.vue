<template>
    <div>
        <h1>Todo Manager</h1>
        <form @submit="onSubmit">
            
            <input type="email" placeholder="Email"/>
            <br>
            <input type="password" placeholder="password">
            <br>
            <input type="submit" value="Register" id=""> -->
        </form>

       <div :key="todo.id" v-for="todo in todos">
           <h3>{{ todo.title }}</h3>
       </div>
    </div>
</template>

<script>
export default {
    name:"Todo",
  
    created(){
        this.todos = [
                {
                title:"Wash clothes",
                isComplete:false,
                timestamp:164323043,
                },
                {
                title:"Read book",
                isComplete:false,
                timestamp:164323043,
                }
            ]
    },

    methods: {
        async onSubmit(e){
            e.preventDefault();
            this.todos.push({
                title:"Read book again",
                isComplete:false,
                timestamp:164323043,
                })
            console.log("submited")
        }
    },
    data(){
    return {
      todos:[
         
      ]
     
    }
  },

}
</script>