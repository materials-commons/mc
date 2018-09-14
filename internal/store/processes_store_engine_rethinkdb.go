package store

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type ProcessesStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewProcessesStoreEngineRethinkdb(session *r.Session) *ProcessesStoreEngineRethinkdb {
	return &ProcessesStoreEngineRethinkdb{Session: session}
}

func (e *ProcessesStoreEngineRethinkdb) AddProcess(process ProcessSchema) (ProcessSchema, error) {
	errMsg := fmt.Sprintf("Unable to add process %+v", process)

	resp, err := r.Table("processes").Insert(process, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return process, err
	}

	var createdProcess ProcessSchema
	err = encoding.Decode(&createdProcess, resp.Changes[0].NewValue)
	return createdProcess, err
}

func (e *ProcessesStoreEngineRethinkdb) GetProcess(processID string) (ProcessExtendedModel, error) {
	var process ProcessExtendedModel
	errMsg := fmt.Sprintf("No such process %s", processID)
	res, err := r.Table("processes").Get(processID).Merge(processDetails).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return process, err
	}

	err = res.One(&process)
	return process, err
}

func processDetails(p r.Term) interface{} {
	return map[string]interface{}{
		"setup": r.Table("process2setup").GetAllByIndex("process_id", p.Field("id")).
			EqJoin("setup_id", r.Table("setups")).Zip().
			Merge(func(row r.Term) interface{} {
				return map[string]interface{}{
					"properties": r.Table("setupproperties").GetAllByIndex("setup_id", row.Field("setup_id")).
						CoerceTo("array"),
				}
			}).CoerceTo("array"),
		"input_samples": r.Table("process2sample"),
	}
}

/*

       input_samples: r.table('process2sample').getAll(process('id'), {index: 'process_id'})
           .filter({'direction': 'in'})
           .eqJoin('sample_id', r.table('samples')).zip()
           .merge(function (sample) {
               return {
                   properties: r.table('propertyset2property')
                       .getAll(sample('property_set_id'), {index: 'property_set_id'})
                       .eqJoin('property_id', r.table('properties')).zip()
                       .orderBy('name')
                       .merge(function (property) {
                           return {
                               best_measure: r.table('best_measure_history')
                                   .getAll(property('best_measure_id'))
                                   .eqJoin('measurement_id', r.table('measurements'))
                                   .zip().coerceTo('array')
                           }
                       }).coerceTo('array'),
                   files: r.table('sample2datafile').getAll(sample('id'), {index: 'sample_id'})
                       .eqJoin('datafile_id', r.table('datafiles')).zip().pluck('id', 'name')
                       .coerceTo('array')
                   // processes: r.table('process2sample').getAll(sample('id'), {index: 'sample_id'})
                   //     .pluck('process_id', 'sample_id').distinct()
                   //     .eqJoin('process_id', r.table('processes')).zip().coerceTo('array')
               }
           }).coerceTo('array'),
       output_samples: r.table('process2sample').getAll(process('id'), {index: 'process_id'})
           .filter({'direction': 'out'})
           .eqJoin('sample_id', r.table('samples')).zip()
           .merge(function (sample) {
               return {
                   properties: r.table('propertyset2property')
                       .getAll(sample('property_set_id'), {index: 'property_set_id'})
                       .eqJoin('property_id', r.table('properties')).zip()
                       .orderBy('name')
                       .merge(function (property) {
                           return {
                               best_measure: r.table('best_measure_history')
                                   .getAll(property('best_measure_id'))
                                   .eqJoin('measurement_id', r.table('measurements'))
                                   .zip().coerceTo('array')
                           }
                       }).coerceTo('array'),
                   files: r.table('sample2datafile').getAll(sample('id'), {index: 'sample_id'})
                       .eqJoin('datafile_id', r.table('datafiles')).zip().pluck('id', 'name')
                       .coerceTo('array')
                   // processes: r.table('process2sample').getAll(sample('id'), {index: 'sample_id'})
                   //     .pluck('process_id', 'sample_id').distinct()
                   //     .eqJoin('process_id', r.table('processes')).zip().coerceTo('array')
               }
           }).coerceTo('array'),
       measurements: r.table('process2measurement').getAll(process('id'), {index: 'process_id'})
           .eqJoin('measurement_id', r.table('measurements')).zip()
           .merge(p2m => {
               return {
                   is_best_measure: r.db('materialscommons').table('best_measure_history')
                       .getAll(p2m('measurement_id'), {index: 'measurement_id'}).count()
               }
           })
           .coerceTo('array'),
       files_count: r.table('process2file').getAll(process('id'), {index: 'process_id'}).count(),
       files: [],
       filesLoaded: false,
       // files: r.table('process2file').getAll(process('id'), {index: 'process_id'})
       //     .eqJoin('datafile_id', r.table('datafiles')).zip()
       //     .merge(f => {
       //         return {
       //             samples: r.table('sample2datafile').getAll(f('id'), {index: 'datafile_id'})
       //                 .eqJoin('sample_id', r.table('samples')).zip()
       //                 .distinct().coerceTo('array')
       //         };
       //     })
       //     .coerceTo('array'),
       input_files: r.table('process2file').getAll(process('id'), {index: 'process_id'})
           .filter({direction: 'in'})
           .eqJoin('datafile_id', r.table('datafiles'))
           .zip().coerceTo('array'),
       output_files: r.table('process2file').getAll(process('id'), {index: 'process_id'})
           .filter({direction: 'out'})
           .eqJoin('datafile_id', r.table('datafiles'))
           .zip().coerceTo('array')
   }
*/
