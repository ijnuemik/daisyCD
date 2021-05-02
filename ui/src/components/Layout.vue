<template>
  <el-row class="layout">
    <el-col class="aside" :span="1">
      <el-menu default-active="2" class="el-menu-vertical" :collapse=true background-color="#545c64" text-color="#fff" active-text-color="#ffd04b">
        <el-submenu index="1">
            <template slot="title">
            <i class="el-icon-location"></i>
            <span slot="title">Navigator One</span>
            </template>
            <el-menu-item-group>
            <span slot="title">Group One</span>
            <el-menu-item index="1-1">item one</el-menu-item>
            <el-menu-item index="1-2">item two</el-menu-item>
            </el-menu-item-group>
            <el-menu-item-group title="Group Two">
            <el-menu-item index="1-3">item three</el-menu-item>
            </el-menu-item-group>
            <el-submenu index="1-4">
            <span slot="title">item four</span>
            <el-menu-item index="1-4-1">item one</el-menu-item>
            </el-submenu>
        </el-submenu>
        <el-menu-item index="2" v-on:click="applications()">
            <i class="el-icon-menu"></i>
            <span slot="title">Manage your applications</span>
        </el-menu-item>
        <el-menu-item index="3" disabled>
            <i class="el-icon-document"></i>
            <span slot="title">Navigator Three</span>
        </el-menu-item>
        <el-menu-item index="4" v-on:click="git()">
            <i class="el-icon-setting"></i>
            <span slot="title">Manage your repositories</span>
        </el-menu-item>
      </el-menu>
    </el-col>
    <el-col class="right" :span="23">
      <el-container class="right-container">
        <el-header class="header">
          <el-menu class="el-menu-horizontal" mode="horizontal">
            <el-menu-item index="1">Processing Center</el-menu-item>
            <el-submenu index="2">
              <template slot="title">Workspace</template>
              <el-menu-item index="2-1">item one</el-menu-item>
              <el-menu-item index="2-2">item two</el-menu-item>
              <el-menu-item index="2-3">item three</el-menu-item>
              <el-submenu index="2-4">
                <template slot="title">item four</template>
                <el-menu-item index="2-4-1">item one</el-menu-item>
                <el-menu-item index="2-4-2">item two</el-menu-item>
                <el-menu-item index="2-4-3">item three</el-menu-item>
              </el-submenu>
            </el-submenu>
            <el-menu-item index="3" disabled>Info</el-menu-item>
            <el-menu-item index="4"><a href="https://www.ele.me" target="_blank">Orders</a></el-menu-item>
          </el-menu>
        </el-header>
        <el-main class="main">
          <el-container class="container">
            <keep-alive>
              <component v-bind:is="currentView"></component>
            </keep-alive>
          </el-container>
        </el-main>
      </el-container>
    </el-col>
  </el-row>
</template>

<style>
  .layout{
    width: 100%;
    height: 100%;
  }
  .aside{
    height: 100%;
  }
  .el-menu-vertical{
    width: 100%;
    height: 100vh;
  }
  .right{
    height: 100vh;
  }
  .right-container{
    height: 100%;
  }
  .header{
    padding: 0px;
  }
  .el-menu-horizontal{
    width: 100%;
    height: 100%;
  }
  .main{
    width: 100%;
    height: 100%;
    padding: 0px;
  }
  .container{
    width: 100%;
    height: 100%;
  }

</style>

<script>
//import Login from '@/components/Login'
//import Settings from '@/components/Settings'
import Git from '@/components/Git'
import Application from '@/components/Application'

export default{
  data() {
    return{
      currentView: '' 
    }
  },
  components: {
    //Login: Login
    //Settings: Settings
  },
  methods: {
    git(){
      this.currentView = Git
    },
    applications(){
      this.currentView = Application
    }
  }
}

</script>