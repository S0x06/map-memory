<template>
  <div class="memory-list-page">
    <el-dialog
      title="记忆详情"
      :visible.sync="memoryDetailBoxVisible"
      width="100%"
      :fullscreen="true"
      >
      <div style="display:flex;flex-direction:column;">
        <div style="display:flex;position:relative;">
            <button v-if="detailMode=='view'" @click="editMemPoint(memDetail.id)" title="编辑" id="edit-mem-point"></button>
            <h3 v-if="detailMode=='view'" style="margin: 0 auto;padding:0px 30px 0px 40px;"><span style="float: left"><img width="40" height="40" :src="'imgs/'+memDetail.icon"/></span>&nbsp;&nbsp;{{memDetail.title}}</h3>
            <div v-if="detailMode=='edit'" style="display:flex;width:100%;">
              <el-select v-model="memDetail.icon" filterable placeholder="记忆图标" @change="setIcon()">
                                 <el-option
                                   v-for="item in available_icons"
                                   :key="item.id"
                                   :label="item.name"
                                   :value="item.id">
                                   <span style="float: left"><img width="20" height="20" :src="'imgs/'+item.id"/></span>
                                   <span style="float: right; color: #8492a6; font-size: 13px">{{ item.name }}</span>
                                 </el-option>
                               </el-select>
              <el-input  v-model="memDetail.title" style="margin: 0 auto;padding:0px 1px 0px 1px;"></el-input>
            </div>

        </div>
        <div style="display:block;margin-top: 10px;">
          <div style="margin-left: auto; font-size:70%;color:gray"><label style="color:#003300;">{{memDetail.nickname}}</label>&nbsp;发布于：{{memDetail.created_at}} &nbsp;<label v-if="autosaveflag">(已自动保存)</label>
          </div>

        </div>
        <div v-if="detailMode=='edit'">

          <wysiwyg @change="autosave" v-model="memDetail.content" />
          <!-- <froala v-model="memDetail.content" :config="option"></froala> -->
          <div class="button-container">
              <el-button style="float: right; margin-right:10px;" type="primary" size="small" @click="saveEdit">保存</el-button>
              <el-button style="float: right; margin-right:10px;" size="small" @click="cancelEdit">取消</el-button>
            </div>
        </div>
        <div v-if="detailMode=='view'" v-html="memDetail.content" style="background-color:#f6f8fa;min-height:300px;overflow:scroll;overflow-x: scroll;font-size:17px;" >
        </div>
      </div>

    </el-dialog>

    <div v-loading="list_loading"  class="center-content">
      <div style="display:flex;margin-top:5px;" @keyup.enter="filterMemory">
           <el-input placeholder="按标题查找"  v-model="title_query"></el-input><el-button @click="filterMemory" type="primary">过滤</el-button>
      </div>
      <el-table
          :border=true
          :data="memorylistdata"
          style="width: 100%">
          <el-table-column
            label="图标"
            prop="title"
            width="50">
            <template slot-scope="scope">
                          <div style="display:flex;">
                            <img width="25" height="25" :src="'/imgs/'+scope.row.icon" onError="this.src='/imgs/question.png';"/>
                          </div>
                        </template>
          </el-table-column>
          <el-table-column
            label="标题"
            prop="title"
            width="170">
          </el-table-column>
          <el-table-column
            label="经度"
            prop="longitude"
            width="100">
          </el-table-column>
          <el-table-column
            label="纬度"
            prop="latitude"
            width="100">
          </el-table-column>
          <el-table-column
            label="创建时间"
            prop="createdAt"
            width="180">
          </el-table-column>
          <el-table-column
            label="操作"
            width="220">
            <template slot-scope="scope">
              <div style="display:flex;">
                <el-button size="small" type="success" @click="go2Map(scope.row.longitude, scope.row.latitude)">前往</el-button>
                <el-button size="small" type="primary" @click="openMemPoint(scope.row.id)">打开</el-button>
                <el-button size="small" type="danger" @click="deleteMemPoint(scope.row.id)">删除</el-button>
              </div>
            </template>
          </el-table-column>

        </el-table>
        <div id="pagi-container">
          <el-pagination style="margin: 0 auto" class="pagination" layout="prev, pager, next" :total="totalCount" :page-size="pageSize"
             v-on:current-change="changePage">
         </el-pagination>
       </div>
    </div>
    <div class="bottom">
      <app-footer></app-footer>
    </div>

  </div>
</template>

<script>
import {AXIOS} from '~/common/http-commons'
import Footer from '~/components/Footer.vue'
import "vue-wysiwyg/dist/vueWysiwyg.css"
import '~/css/map.public.css'
import state from '~/common/state'
import base from '~/mixins/base'
import _ from 'lodash'
import Qs from 'qs'

export default {
    name: 'MemoryList',
    components: {
      'app-footer': Footer
    },
    mixins: [base],
    computed: {
    option() {
        return {
          imageUpload: false,
          heightMin:350,
          saveInterval: 2500,
          saveMethod: 'POST',
          saveParam: 'memDetail.content',
          saveParams:{
            'id': this.memDetail.id
          },
          requestWithCredentials: true,
          requestWithCORS: true,
          saveURL: this.base_service_url+'/api/v1/memory-content'
        };
    }
    },
    data() {
      return {
        base_service_url:'',
        title_query:'',
        state,
        memoryDetailBoxVisible: false,
        detailMode:'view',
        totalCount: 0,
        currentPage: 1,
        pageSize: 0,
        memorylistdata:[],
        admin_request_id: '',
        memDetail:{},
        list_loading: false,
        autosaveflag:false
      }
    },
    methods: {
      autosave:_.debounce(function () {
          var data = {'id': this.memDetail.id, 'content':this.memDetail.content};
          AXIOS.put('/api/v1/memory-content'
          ,Qs.stringify(data), 
          {headers:{'Content-Type':'application/x-www-form-urlencoded'}}).then(response=>{
            if(response.data.ok==true){
              this.autosaveflag = true;
              var count = 2;
              var storeThis = this;
              var refreshIntervalId = setInterval(function() {
              count = count - 1;
              if(count==0){
                clearInterval(refreshIntervalId);
                storeThis.autosaveflag = false;
              }
             }, 1000);
            }
          }).catch(e=>{

          })
      }, 2500),
      filterMemory(){
        this.getMemoryListData();
      },
      saveEdit(){
        // update with: id, title, content, that's it
        AXIOS.put('/api/v1/memory/'+this.memDetail.id,{
          title: this.memDetail.title,
          content: this.memDetail.content,
          icon: this.memDetail.icon
        }).then(response=>{
          if(response.data.ok==true){
            this.$notify({
                            title: '成功',
                            type: 'success',
                            message: response.data.message
                          });
            this.cancelEdit();

          }else{
            this.$notify.error({
              title: '错误',
              message: response.data.message
            })
          }
        }).catch(e=>{
          this.$notify.error({
            title: '错误',
            message: '未知错误'
          })
        })
      },
      cancelEdit(){
        this.detailMode = 'view';
        this.loadMemoryDetailById(this.memDetail.id);
      },
      editMemPoint(id){
        // make title and content edit mode
        this.detailMode = 'edit';
      },
      loadMemoryDetailById(id){
        AXIOS.get('/api/v1/memory/'+id).then(response=>{
          if(response.data.ok==true){
            this.memDetail = response.data.data;
          }else{
            this.$notify.error({
              title: '错误',
              message: response.data.message
            })
          }
        }).catch(e=>{
          this.$notify.error({
            title: '错误',
            message: '未知错误'
          })
        })
      },
      changePage: function (currentPage) {
       this.currentPage = currentPage;
       this.getMemoryListData();
      },
      go2Map(longitude, latitude){
        this.state.longitude = longitude;
        this.state.latitude = latitude;
        this.$router.push('/working');
      },
      openMemPoint(id){
        this.memoryDetailBoxVisible = true;
        this.loadMemoryDetailById(id);
      },
      deleteMemPoint(id){
        this.$confirm('确认删除该记忆点？')
                    .then(_ => {
                      AXIOS.delete('/api/v1/memory/'+id).then(response=>{
                        if(response.data.ok==true){
                          this.$notify({
                                          title: '成功',
                                          type: 'success',
                                          message: response.data.message
                                        });

                          this.getMemoryListData();
                        }else{
                          this.$notify.error({
                            title: '错误',
                            message: response.data.message
                          })
                        }
                      }).catch(e=>{
                        this.$notify.error({
                          title: '错误',
                          message: '未知错误'
                        })
                      })
                      done();
                    }).catch(_ => {

                    });
      },
      getMemoryListData(){
        var params = {};
        params.page = this.currentPage;
        params.title_query = this.title_query;
        this.list_loading=true;
        AXIOS.get('/api/v1/memory-my',{
            params: params
        }).then(response=>{
          if(response.data.ok==true){
            this.memorylistdata = response.data.data.data;
            this.totalCount = Number(response.data.data.totalCount);
            this.currentPage = Number(response.data.data.currentPage);
            this.pageSize = Number(response.data.data.pageSize);
          }else{
            this.$notify.error({
              title: '错误',
              message: response.data.message
            });
          }
          this.list_loading=false;

        }).catch(e=>{
          this.$notify.error({
            title: '错误',
            message: '未知错误'
          });
          this.list_loading=false;
        })
      }
    },
    created () {

    },
    mounted () {
      this.checklogin();
      this.getBaseServiceUrl();
      this.getMemoryListData();
    }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.memory-list-page{
  flex-flow: column;
  width: 100%;
  display: flex;
  left: 0;
  position:absolute;
  background-color: #f1f4f5;
}


.bottom{
  flex: 1;
  background-color: #F3F0EC;
}
.center-content{
  margin: 0 auto;
  background-color: white;
  width: 800px;
  min-height: 700px;
}

</style>
